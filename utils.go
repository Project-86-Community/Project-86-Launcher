package p86l

import (
	"runtime"
	"strings"

	webview "github.com/crgimenes/glaze"
	_ "github.com/crgimenes/glaze/embedded"
)

// WebviewRequest is sent to the dedicated webview goroutine.
type WebviewRequest struct {
	Title  string
	Source string
	Reply  chan string
}

// RunWebviewThread starts a dedicated OS-locked goroutine that serves
// webview requests. Call this once before guigui.Run().
// Send requests to the returned channel to open a webview window.
func RunWebviewThread() chan<- WebviewRequest {
	ch := make(chan WebviewRequest)
	go func() {
		// Lock this goroutine to its OS thread permanently.
		runtime.LockOSThread()
		for req := range ch {
			req.Reply <- webviewOpen(req.Title, req.Source)
		}
	}()
	return ch
}

func webviewOpen(title, source string) string {
	selectedURL := make(chan string, 1)

	w, err := webview.New(true)
	if err != nil {
		return ""
	}
	defer func() {
		w.Unbind("reportURL")
		w.Destroy()
	}()

	w.SetTitle(title)
	w.SetSize(1024, 768, webview.HintNone)

	if err := w.Bind("reportURL", func(url string) {
		if strings.Contains(url, "/releases/download/") {
			select {
			case selectedURL <- url:
				w.Terminate()
			default:
			}
		}
	}); err != nil {
		return ""
	}

	w.Init(`
		(function() {
			function setup() {
				document.addEventListener('mousedown', function(e) {
					var a = e.target.closest('a[href]');
					if (a && a.href && window.reportURL) {
						window.reportURL(a.href);
					}
				}, true);

				var lastURL = window.location.href;
				setInterval(function() {
					var cur = window.location.href;
					if (cur !== lastURL) {
						lastURL = cur;
						if (window.reportURL) window.reportURL(cur);
					}
				}, 300);
			}

			if (document.readyState === 'loading') {
				document.addEventListener('DOMContentLoaded', setup);
			} else {
				setup();
			}
		})();
	`)

	w.Navigate(source)
	w.Run()

	select {
	case url := <-selectedURL:
		return url
	default:
		return ""
	}
}
