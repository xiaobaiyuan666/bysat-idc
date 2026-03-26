package server

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerCompiledWeb(engine *gin.Engine, adminDir string, portalDir string) {
	adminMounted := hasCompiledIndex(adminDir)
	portalMounted := hasCompiledIndex(portalDir)

	if adminMounted {
		engine.GET("/admin", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/admin/")
		})
	}

	if portalMounted {
		engine.GET("/portal", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/portal/")
		})
	}

	if adminMounted {
		engine.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/admin/")
		})
	} else if portalMounted {
		engine.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/portal/")
		})
	}

	if !adminMounted && !portalMounted {
		return
	}

	engine.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Status(http.StatusNotFound)
			return
		}

		requestPath := c.Request.URL.Path
		switch {
		case adminMounted && strings.HasPrefix(requestPath, "/admin/"):
			serveCompiledSPA(c, adminDir, strings.TrimPrefix(requestPath, "/admin/"))
		case portalMounted && strings.HasPrefix(requestPath, "/portal/"):
			serveCompiledSPA(c, portalDir, strings.TrimPrefix(requestPath, "/portal/"))
		default:
			c.Status(http.StatusNotFound)
		}
	})
}

func hasCompiledIndex(dir string) bool {
	if strings.TrimSpace(dir) == "" {
		return false
	}

	_, err := os.Stat(filepath.Join(dir, "index.html"))
	return err == nil
}

func serveCompiledSPA(c *gin.Context, dir string, requested string) {
	indexPath := filepath.Join(dir, "index.html")
	cleanRelative := strings.TrimPrefix(path.Clean("/"+requested), "/")
	if cleanRelative == "" || cleanRelative == "." {
		c.File(indexPath)
		return
	}

	candidate := filepath.Join(dir, filepath.FromSlash(cleanRelative))
	relativeToBase, err := filepath.Rel(filepath.Clean(dir), filepath.Clean(candidate))
	if err != nil || relativeToBase == ".." || strings.HasPrefix(relativeToBase, ".."+string(os.PathSeparator)) {
		c.Status(http.StatusNotFound)
		return
	}

	info, statErr := os.Stat(candidate)
	if statErr == nil && !info.IsDir() {
		c.File(candidate)
		return
	}

	if path.Ext(cleanRelative) != "" {
		c.Status(http.StatusNotFound)
		return
	}

	c.File(indexPath)
}
