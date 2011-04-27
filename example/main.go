// example of HTTP server that uses the captcha package.
package main

import (
	"fmt"
	"github.com/dchest/captcha"
	"http"
	"io"
	"log"
	"template"
)

var formTemplate = template.MustParse(formTemplateSrc, nil)

func showFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	d := struct {
		CaptchaId  string
		JavaScript string
	}{
		captcha.New(),
		formJavaScript,
	}
	if err := formTemplate.Execute(w, &d); err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
	}
}

func processFormHandler(w http.ResponseWriter, r *http.Request) {
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		io.WriteString(w, "Wrong captcha solution! No robots allowed!\n")
	} else {
		io.WriteString(w, "Great job, human! You solved the captcha.\n")
	}
	io.WriteString(w, "<br><a href='/'>Try another one</a>")
}

func main() {
	http.HandleFunc("/", showFormHandler)
	http.HandleFunc("/process", processFormHandler)
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	fmt.Println("Server is at localhost:8666")
	if err := http.ListenAndServe(":8666", nil); err != nil {
		log.Fatal(err)
	}
}

const formJavaScript = `
function playAudio() {
	var e = document.getElementById('audio')
	e.style.display = 'block';
	e.play();
	return false;
}

function reload() {
	function setSrcQuery(e, q) {
		var src  = e.src;
		var p = src.indexOf('?');
		if (p >= 0) {
			src = src.substr(0, p);
		}
		e.src = src + "?" + q
	}
	setSrcQuery(document.getElementById('image'), "reload=" + (new Date()).getTime());
	setSrcQuery(document.getElementById('audio'), (new Date()).getTime());
	return false;
}
`

const formTemplateSrc = `<!doctype html>
<head><title>Captcha Example</title></head>
<body>
<script>
{JavaScript}
</script>
<form action="/process" method=post>
<p>Type the numbers you see in the picture below:</p>
<p><img id=image src="/captcha/{CaptchaId}.png" alt="Captcha image"></p>
<a href="#" onclick="reload()">Reload</a> | <a href="#" onclick="playAudio()">Play Audio</a>
<audio id=audio controls style="display:none" src="/captcha/{CaptchaId}.wav" preload=none>
  You browser doesn't support audio.
  <a href="/captcha/download/{CaptchaId}.wav">Download file</a> to play it in the external player.
</audio>
<input type=hidden name=captchaId value="{CaptchaId}"><br>
<input name=captchaSolution>
<input type=submit value=Submit>
</form>
`
