package server

import (

	"embed"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"mini/internal/config"
	"mini/internal/files"

)

//go:embed ../../web/templates/*.tmpl
var tmplFS embed.FS

// Handlers groups deps for HTTP handlers
type Handlers struct {
	cfg *config.App
	tmpl *template.Template
}

func NewHandlers(cfg *config.App) *Handlers {
	t := template.Must(template.ParseFS(tmplFS, "web/templates/*.tmpl"))
	return &Handlers {
		cfg:  cfg,
		tmpl: t,
	}
}

// RedirectRoot sends "/" to main listing
func (h *Handlers) RedirectRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view?dir=/", http.StatusFound)
}

// ViewDir list directory contents
func (h *Handlers) ViewDir(w http.ResponseWriter, r *http.Request) {
	dir := filepath.Clean(r.URL.Query().Get("dir")
	if dir == "" {
		dir = "/"
	}
	parent := filepath.Dir(dir)
	if parent == "." {
		parent = "/"
	}

	abs := filepath.Clean(filepath.Join(h.cfg.StorageDir, dir))
	lst, err := files.List(abs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = h.tmpl.ExecuteTemplate(w, "dirlist.html.tmpl", map[string]any{
		"Dir": dir, "Parent": parent, "Files": lst,
		"Title": "Directory Listing for " + dir,
	})
}

// GetFile serves a file to the client with "Content-Disposition attachment"
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Disposition

func (h *Handlers) GetFile(w http.ResponseWriter, r *http.Request) {
	f := filepath.Clean(r.URL.Query().Get("file"))
	abs := filepath.Clean(filepath.Join(h.cfg.StorageDir, f))
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(abs))
	http.ServeFile(w, r, abs)
}

// Upload saves multipart file
// multipart/form-data info
// https://www.w3.org/TR/html401/interact/forms.html#h-17.13.4

func (h *Handlers) Upload (w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultiPartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	dir := filepath.Clean(r.FormValue("directory"))
	for _, fh := range r.MultipartForm.File["file-upload"] {
		src, _ = fh.Open()
		dst := filepath.Clean(filepath.Join(h.cfg.StorageDir, fh.Filename))

		if err := files.SaveUpload(dst, src); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/view?dir="+dir, http.StatusFound)
	}
}

// Delete - delete, redirect to root
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	dir := filepath.Clean(r.FormValue("directory"))
	fn := filepath.Clean(r.FormValue("filename"))
	abs := filepath.Clean(filepath.Join(h.cfg.StorageDir, dir, fn))
	_ = os.Remove(abs)
	http.Redirect(w, r, "/view?dir="+dir, http.StatusFound)
}

// compareHash wraps bcrypt password check 
func compareHash(hash, pw string) bool {
	return auth.ComparePassword(hash, pw)
}