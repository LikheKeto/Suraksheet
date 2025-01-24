package document

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/LikheKeto/Suraksheet/service/auth"
	"github.com/LikheKeto/Suraksheet/types"
	"github.com/LikheKeto/Suraksheet/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	store     types.DocumentStore
	userStore types.UserStore
	binStore  types.BinStore
	minio     *minio.Client
	rmqChan   *amqp.Channel
	rmq       amqp.Queue
}

func NewHandler(documentStore types.DocumentStore,
	userStore types.UserStore, binStore types.BinStore,
	minio *minio.Client, rmqChan *amqp.Channel, rmq amqp.Queue) *Handler {
	return &Handler{
		store:     documentStore,
		userStore: userStore,
		binStore:  binStore,
		minio:     minio,
		rmqChan:   rmqChan,
		rmq:       rmq,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.MethodFunc(http.MethodGet, "/document/{documentID}/asset", auth.WithJWTAuth(h.handleGetImage, h.userStore))
	router.MethodFunc(http.MethodGet, "/document/{documentID}", auth.WithJWTAuth(h.handleGetDocument, h.userStore))
	router.MethodFunc(http.MethodPost, "/document", auth.WithJWTAuth(h.handleInsertDocument, h.userStore))
	router.MethodFunc(http.MethodPatch, "/document", auth.WithJWTAuth(h.handleEditDocument, h.userStore))
	router.MethodFunc(http.MethodDelete, "/document", auth.WithJWTAuth(h.handleDeleteDocument, h.userStore))
}

func (h *Handler) handleGetImage(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	documentIDStr := chi.URLParam(r, "documentID")
	documentID, err := strconv.Atoi(documentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid document id"))
		return
	}
	owner, _ := h.store.GetDocumentOwner(documentID)
	if owner != user.ID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("document does not belong to user"))
		return
	}
	document, err := h.store.GetDocumentByID(documentID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	objectName := path.Join(utils.HashString(user.Email), strconv.Itoa(document.BinID), utils.HashString(document.ReferenceName))
	obj, err := utils.GetObject(r.Context(), h.minio, objectName)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	stat, err := obj.Stat()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Errorf("unable to get object stats: %w", err))
		return
	}
	defer obj.Close()
	if _, err := io.Copy(w, obj); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to write object to response: %w", err))
		return
	}
	w.Header().Set("Content-Type", stat.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", stat.Key))
}

func (h *Handler) handleGetDocument(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	documentIDStr := chi.URLParam(r, "documentID")
	documentID, err := strconv.Atoi(documentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid document id"))
		return
	}
	owner, _ := h.store.GetDocumentOwner(documentID)
	if owner != user.ID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("document does not belong to user"))
		return
	}
	document, err := h.store.GetDocumentByID(documentID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, document)
}

func (h *Handler) handleInsertDocument(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	r.ParseMultipartForm(10 << 20)
	referenceName := r.Form.Get("referenceName")
	binIDStr := r.Form.Get("binID")

	binID, err := strconv.Atoi(binIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid bin"))
		return
	}

	if err := h.store.ReferenceNameExistsInBin(referenceName, binID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get file from request: %v", err), http.StatusInternalServerError)
		return
	}

	// Upload file to MinIO
	fileKey := path.Join(utils.HashString(user.Email), strconv.Itoa(binID), utils.HashString(referenceName))
	err = utils.UploadToMinio(r.Context(), h.minio, file, fileHeader, fileKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	docID, err := h.store.InsertDocument(types.Document{
		BinID:         binID,
		Url:           "",
		Name:          fileHeader.Filename,
		ReferenceName: referenceName,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("unable to insert document: %v", err))
		return
	}
	// queue for extraction
	parts := strings.Split(fileHeader.Filename, ".")
	var extension string
	if len(parts) > 1 {
		extension = parts[len(parts)-1]
	}

	err = utils.QueueForExtraction(h.rmqChan, h.rmq, docID, fileKey, extension)
	if err != nil {
		utils.WriteJSON(w, http.StatusCreated, fmt.Errorf("unable to queue for extraction: %v", err))
	}

	doc := types.Document{
		ID:            int(docID),
		Name:          fileHeader.Filename,
		ReferenceName: referenceName,
		BinID:         binID,
		Url:           "",
		Extract:       "",
		CreatedAt:     time.Now(),
	}
	utils.WriteJSON(w, http.StatusCreated, doc)
}

func (h *Handler) handleEditDocument(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}
	var payload types.EditDocumentPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	doc, err := h.store.GetDocumentByID(payload.Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if doc.ReferenceName == payload.ReferenceName {
		utils.WriteJSON(w, http.StatusNoContent, nil)
		return
	}

	bin, err := h.binStore.GetBinById(doc.BinID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if bin.OwnerID != user.ID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("document doesn't belong to user"))
		return
	}

	err = h.store.UpdateDocumentName(payload.Id, payload.ReferenceName)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	old := path.Join(utils.HashString(user.Email), strconv.Itoa(bin.ID), utils.HashString(doc.ReferenceName))
	new := path.Join(utils.HashString(user.Email), strconv.Itoa(bin.ID), utils.HashString(payload.ReferenceName))
	err = utils.RenameObject(r.Context(), h.minio, old, new)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("unable to rename document: %v", err))
		err = h.store.UpdateDocumentName(payload.Id, doc.ReferenceName)
		if err != nil {
			log.Printf("unable to revert file name in db: %v\n", err)
		}
		return
	}
	err = utils.DeleteObject(r.Context(), h.minio, old)
	if err != nil {
		log.Printf("unable to delete file in minio: %v\n", err)
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleDeleteDocument(w http.ResponseWriter, r *http.Request) {
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
	doc, err := h.store.GetDocumentByID(payload.Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	bin, err := h.binStore.GetBinById(doc.BinID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if bin.OwnerID != user.ID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("document doesn't belong to user"))
		return
	}
	err = utils.DeleteObject(r.Context(), h.minio, path.Join(utils.HashString(user.Email),
		strconv.Itoa(doc.BinID), utils.HashString(doc.ReferenceName)))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("unable to delete object: %v", err))
		return
	}
	err = h.store.DeleteDocumentByID(payload.Id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}
