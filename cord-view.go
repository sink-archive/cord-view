package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/webview/webview"
)

const JS_OVERRIDES string = `
window.CordView = {
  native: { fakefetch: temp_fakefetch },
  fetch: (url) =>
    CordView.native
      .fakefetch(url)
      .then(JSON.parse)
      .then((res) => ({
        text: () => Promise.resolve(res.text),
        json: () => Promise.resolve(JSON.parse(res.text)),
        status: res.status,
      })) /* 
  getFakeWindow: () => {
    const fakeWindow = Object.assign({}, window);
    fakeWindow.window = fakeWindow.self = fakeWindow;
    Object.assign(fakeWindow, CordView);
    return fakeWindow;
  }, */,
  eval: function (code) {
    return eval.call(
      this,
      "((eval, fetch)=>{return " +
        code
          .replaceAll("=window.eval", "=CordView.eval")
          .replaceAll("window.eval(", "CordView.eval(") +
        "\n})(CordView.eval, CordView.fetch)"
    );
  },
};

delete window.temp_fakefetch;
`

const CC_URL string = "https://raw.githubusercontent.com/Cumcord/builds/main/build.js"

func fakeFetch(url string) string {
	print("\n\n\u001b[31mfake fetch called for: " + url + "\n\n\u001b[37m")
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	txt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	ret := struct {
		Text   string `json:"text"`
		Status int `json:"status"`
	}{
		Text:   string(txt),
		Status: resp.StatusCode,
	}

	fmt.Printf("\n\n\u001b[32m%d\n\n\u001b[37m", ret.Status)

	json, err := json.Marshal(ret)
	if err != nil {
		return ""
	}
	return string(json)
}

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("cord-view")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://discord.com/channels/@me")

	// inject initial fakefetch, but this will be overwritten later
	w.Bind("temp_fakefetch", fakeFetch)

	go (func() {
		time.Sleep(3 * time.Second)
		w.Dispatch(func() {
			// finish off fake fetch, and inject fakeeval
			w.Eval(JS_OVERRIDES)
			w.Eval("CordView.fetch('" + CC_URL + "').then(f => f.text()).then(CordView.eval)")
		})
	})()

	w.Run()
}
