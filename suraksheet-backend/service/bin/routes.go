package bin

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/LikheKeto/Suraksheet/service/auth"
	"github.com/LikheKeto/Suraksheet/types"
	"github.com/LikheKeto/Suraksheet/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
)

type Handler struct {
	store         types.BinStore
	userStore     types.UserStore
	documentStore types.DocumentStore
	minio         *minio.Client
}

func NewHandler(store types.BinStore, userStore types.UserStore, documentStore types.DocumentStore, minio *minio.Client) *Handler {
	return &Handler{store: store, userStore: userStore, documentStore: documentStore, minio: minio}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.MethodFunc(http.MethodGet, "/bins/{binID}", auth.WithJWTAuth(h.handleGetDocumentsInBin, h.userStore))
	router.MethodFunc(http.MethodGet, "/bins", auth.WithJWTAuth(h.handleGetBins, h.userStore))
	router.MethodFunc(http.MethodPost, "/bins", auth.WithJWTAuth(h.handleCreateBin, h.userStore))
	router.MethodFunc(http.MethodPatch, "/bins", auth.WithJWTAuth(h.handleEditBin, h.userStore))
	router.MethodFunc(http.MethodDelete, "/bins", auth.WithJWTAuth(h.handleDeleteBin, h.userStore))
}

func (h *Handler) handleGetDocumentsInBin(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	binIDStr := chi.URLParam(r, "binID")
	BinID, err := strconv.Atoi(binIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid bin %s", binIDStr))
		return
	}
	bin, err := h.store.GetBinById(BinID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if bin.OwnerID != user.ID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("bin doesn't belong to user"))
		return
	}
	documents, err := h.documentStore.GetDocumentsInBin(BinID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, documents)
}

func (h *Handler) handleGetBins(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	bins, err := h.store.GetBinsByUser(user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, bins)
}

func (h *Handler) handleCreateBin(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	var payload types.CreateBinPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}
	err = h.store.CreateBin(payload.Name, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleEditBin(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	var payload types.EditBinPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}
	err = h.store.UpdateBinName(payload.Id, user.ID, payload.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("unable to update name: %v", err))
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleDeleteBin(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	var payload types.DeleteBinDocPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err = utils.DeleteDir(r.Context(), h.minio, path.Join(utils.HashString(user.Email), strconv.Itoa(payload.Id)))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("unable to delete bin: %v", err))
	}

	err = h.store.DeleteBin(payload.Id, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("unable to delete bin: %v", err))
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}
