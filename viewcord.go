package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/webview/webview"
)

func getCcBundle() string {
	resp, err := http.Get("https://raw.githubusercontent.com/Cumcord/builds/main/build.js")
	if (err != nil) {return "(() => {})()"}

	txt, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {return "(() => {})()"}
	return string(txt)
}

func overrideFetch(bundle string) string {
	return "((fetch) => {" + bundle + "\n})(fake_fetch)"
}

const JS_FAKEFETCH_OVERRIDE string = "{const a=window.fake_fetch;window.fake_fetch=async b=>{const c=await a(b);return{text:()=>Promise.resolve(c),json:()=>Promise.resolve(JSON.parse(c)),status:200}}}"

func fakeFetch(url string) string {
	print("\n\nfake fetch called for: " + url + "\n\n")
	resp, err := http.Get(url)
	if (err != nil) {return ""}

	txt, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {return ""}

	return string(txt)
}

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Viewcord")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://discord.com/channels/@me")
	
	w.Bind("fake_fetch", fakeFetch)

	go (func () {
		ccBundle := overrideFetch(getCcBundle())
		time.Sleep(3 * time.Second)
		w.Dispatch(func() {
			w.Eval(JS_FAKEFETCH_OVERRIDE)
			w.Eval(ccBundle)
		})
	})()

	w.Run()
}
