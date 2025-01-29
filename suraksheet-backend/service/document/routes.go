package document

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/LikheKeto/Suraksheet/service/auth"
	"github.com/LikheKeto/Suraksheet/types"
	"github.com/LikheKeto/Suraksheet/utils"
	"github.com/elastic/go-elasticsearch/v8"
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
	esClient  *elasticsearch.Client
}

func NewHandler(documentStore types.DocumentStore,
	userStore types.UserStore, binStore types.BinStore,
	minio *minio.Client, rmqChan *amqp.Channel, rmq amqp.Queue, esClient *elasticsearch.Client) *Handler {
	return &Handler{
		store:     documentStore,
		userStore: userStore,
		binStore:  binStore,
		minio:     minio,
		rmqChan:   rmqChan,
		rmq:       rmq,
		esClient:  esClient,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.MethodFunc(http.MethodGet, "/document/{documentID}/asset", auth.WithJWTAuth(h.handleGetImage, h.userStore))
	router.MethodFunc(http.MethodGet, "/document/{documentID}", auth.WithJWTAuth(h.handleGetDocument, h.userStore))
	router.MethodFunc(http.MethodPost, "/document", auth.WithJWTAuth(h.handleInsertDocument, h.userStore))
	router.MethodFunc(http.MethodPatch, "/document", auth.WithJWTAuth(h.handleEditDocument, h.userStore))
	router.MethodFunc(http.MethodDelete, "/document", auth.WithJWTAuth(h.handleDeleteDocument, h.userStore))
	router.MethodFunc(http.MethodGet, "/document/search", auth.WithJWTAuth(h.handleSearchDocuments, h.userStore))
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
	language := r.Form.Get("language")

	if !(language == "eng" || language == "nep") {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("language not supported"))
		return
	}

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

	document, err := h.store.InsertDocument(types.Document{
		BinID:         binID,
		Url:           "",
		Name:          fileHeader.Filename,
		ReferenceName: referenceName,
		Language:      language,
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

	err = utils.QueueForExtraction(h.rmqChan, h.rmq, utils.ExtractionArgs{
		DocID:     document.ID,
		UserID:    user.ID,
		FileKey:   fileKey,
		Extension: extension,
		Language:  language,
	})
	if err != nil {
		utils.WriteJSON(w, http.StatusCreated, fmt.Errorf("unable to queue for extraction: %v", err))
	}
	utils.WriteJSON(w, http.StatusCreated, document)
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

func (h *Handler) handleSearchDocuments(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ExtractUserFromContext(r)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("query parameter 'q' is required"))
		return
	}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"match": map[string]interface{}{"text": query}},
				},
				"filter": []map[string]interface{}{
					{"term": map[string]interface{}{"user_id": user.ID}},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to encode search query: %v", err))
		return
	}

	res, err := h.esClient.Search(
		h.esClient.Search.WithContext(context.Background()),
		h.esClient.Search.WithIndex("documents"),
		h.esClient.Search.WithBody(&buf),
		h.esClient.Search.WithPretty(),
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("search query failed: %v", err))
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		utils.WriteError(w, res.StatusCode, fmt.Errorf("error from elasticsearch"))
		return
	}

	var esRes map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract document IDs from Elasticsearch response
	var docIDs []int
	if hits, ok := esRes["hits"].(map[string]interface{})["hits"].([]interface{}); ok {
		for i, hit := range hits {
			if i > 3 {
				break
			}
			if hitMap, ok := hit.(map[string]interface{}); ok {
				if source, ok := hitMap["_source"].(map[string]interface{}); ok {
					if docID, ok := source["document_id"].(float64); ok {
						docIDs = append(docIDs, int(docID))
					}
				}
			}
		}
	}

	if len(docIDs) == 0 {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	documents, err := h.store.FetchDocumentsFromDB(docIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(documents)
}
