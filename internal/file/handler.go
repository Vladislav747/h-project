package file

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type FileHandler struct {
	Logger      *slog.Logger
	FileService Service
}

func NewFileHandler(service Service, logger *slog.Logger) *FileHandler {
	return &FileHandler{
		FileService: service,
		Logger:      logger,
	}
}

func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GetFile")

	h.Logger.Debug("get note_uuid from URL")

	fileId := r.URL.Query().Get("fileId")
	if fileId == "" {
		return errors.New("no file id provided")
	}

	h.Logger.Debug("get field from context")

	f, err := h.FileService.GetFile(r.Context(), BUCKET_NAME, fileId)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f.Name))
	//Отдаем с таким же заголовком которые запрашивали
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	w.Write(f.Bytes)

	return nil
}

func (h *FileHandler) CreateFile(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CreateFile")
	w.Header().Set("Content-Type", "form/json")

	// TODO maximum file size
	/**
	- r.ParseMultipartForm() - это метод из стандартной библиотеки Go,
	который разбирает данные, отправленные в HTTP-запросе с типом multipart/form-data.
	- 32 << 20 - это способ задать максимальный размер памяти, который может быть использован для хранения данных из формы. В данном случае это 32 мегабайта (32 * 2^20 = 33554432 байт). Это нужно, чтобы предотвратить атаки с большими объемами данных, которые могут исчерпать память сервера.
	*/
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}
	h.Logger.Debug("decode create upload fileInfo dto")

	//Если не пришел заголовок с файлом
	files, ok := r.MultipartForm.File["file"]
	if !ok || len(files) == 0 {
		return errors.New("file required")
	}

	fileInfo := files[0]
	fileReader, err := fileInfo.Open()
	dto := CreateFileDTO{
		Name:   fileInfo.Filename,
		Size:   fileInfo.Size,
		Reader: fileReader,
	}

	err = h.FileService.Create(r.Context(), r.Form.Get("node_uuid"), dto)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DeleteFile")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get fileId from context")

	fileId := r.URL.Query().Get("note_uuid")
	if fileId == "" {
		return errors.New("file required")
	}

	err := h.FileService.Delete(r.Context(), BUCKET_NAME, fileId)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
