package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/webview/webview"
)

/*

{
  const oldFakeFetch = window.fake_fetch;
  window.fake_fetch = async (url) => {
    const txt = await oldFakeFetch(url);
    return {
      text: () => Promise.resolve(txt),
      json: () => Promise.resolve(JSON.parse(txt)),
      status: 200,
    };
  };
  const oldEval = eval;
  window.eval = function (code) {
    return oldEval.call(
      this,
      `((fetch,eval)=>{${code}
})(fetch,eval)`
    );
  };
}

*/
const JS_OVERRIDES string = "const a=fakefetch,b=eval;window.fakefetch=b=>a(b).then(a=>({text:()=>Promise.resolve(a),json:()=>Promise.resolve(JSON.parse(a)),status:200})),window.fakeeval=function(a){return b.call(this,`((fetch,eval)=>{${a}\n})(fakefetch,fakeeval)`)}"

const CC_URL string = "https://raw.githubusercontent.com/Cumcord/builds/main/build.js"

func fakeFetch(url string) string {
	print("\n\nfake fetch called for: " + url + "\n\n")
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	txt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(txt)
}

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("cord-view")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://discord.com/channels/@me")

	// inject initial fakefetch, but this will be overwritten later
	w.Bind("fakefetch", fakeFetch)

	go (func() {
		time.Sleep(3 * time.Second)
		w.Dispatch(func() {
			// finish off fake fetch, and inject fakeeval
			w.Eval(JS_OVERRIDES)
			w.Eval("fakefetch('" + CC_URL + "').then(f => f.text()).then(fakeeval)")
		})
	})()

	w.Run()
}
