package monitor

import (
	"log"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var (
	html string = `    <!DOCTYPE html>
<html>
<head>
    <title>Kafdrop: Consumer: jaeger-injest-1</title>
    <link type="text/css" rel="stylesheet" href="/kafdrop/css/bootstrap.min.css"/>
    <link type="text/css" rel="stylesheet" href="/kafdrop/css/font-awesome.min.css"/>
    <link type="text/css" rel="stylesheet" href="/kafdrop/css/global.css"/>

    <script src="/kafdrop/js/jquery.min.js"></script>
    <script src="/kafdrop/js/popper.min.js"></script>
    <script src="/kafdrop/js/bootstrap.min.js"></script>
    <script src="/kafdrop/js/global.js"></script>
    <script async defer src="/kafdrop/js/github-buttons.js"></script>

</head>
<body>

<div class="pb-2 mt-5 mb-4 border-bottom border-secondary">
    <div class="container">
        <div class="container-fluid pl-0">
            <div id="header-title-line" class="row">
                <div id="logo" class="col-md-1">
                    <img alt="logo" height="100%" src="/kafdrop/images/kafdrop-logo.svg"/>
                </div>
                <div id="title" class="col-md-10">
                    <h1 class="app-name brand mb-0">
                        <a href="/kafdrop/">Kafdrop</a>
                    </h1>
                </div>
                <div id="github-star" class="col-md-1">
                    <a class="github-button" href="https://github.com/obsidiandynamics/kafdrop" data-show-count="false"
                       aria-label="Star Kafdrop on GitHub" data-color-scheme="dark">Star</a>
                </div>
                <script>
                    $(document).ready(function(){
                        setTimeout(function() { restyle(0); });
                    });

                    function restyle(retries) {
                        var githubStarSpan = document.querySelector('#github-star span');
                        if (githubStarSpan != null) {
                            var shadowRoot = githubStarSpan.shadowRoot;
                            shadowRoot.querySelector('.btn')
                                .setAttribute('style', 'color:#00f0fe; background-image:none; background-color:#222; border-color:#222');
                            shadowRoot.querySelector('.social-count')
                                .setAttribute('style', 'color:#222; background-image:none; background-color:#00f0fe; border-color:#00f0fe');
                            shadowRoot.querySelector('.social-count b')
                                .setAttribute('style', 'border-right-color:#00f0fe');
                            shadowRoot.querySelector('.social-count i')
                                .setAttribute('style', 'border-right-color:#00f0fe');
                        } else {
                            setTimeout(function() { restyle(retries + 1); }, retries * 10);
                        }
                    }
                </script>
            </div>
        </div>
    </div>
</div>
<div class="container">
<h2>Kafka Consumer: jaeger-injest-1</h2>
<div class="container-fluid pl-0">
    <div id="overview">
        <h3>Overview</h3>
        <table class="table table-bordered overview">
            <tbody>
            <tr>
                <td>Topics</td>
                <td>1</td>
            </tr>
            </tbody>
        </table>
    </div>
    <div id="topics">
            <h3>    <a href="#topic-0-table" class="toggle-link" data-toggle-target="#topic-0-table"><i
                class="fa fa-chevron-circle-down"></i></a>
 Topic: <a href="/kafdrop/topic/jaeger-spans">jaeger-spans</a></h3>
            <div id="topic-0-table">
                <p>
                <table class="table table-bordered table-sm">
                    <thead>
                    <tr>
                        <th>Partition</th>
                        <th>First Offset</th>
                        <th>Last Offset</th>
                        <th>Consumer Offset</th>
                        <th>Lag</th>
                    </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>0</td>
                            <td>3947522053</td>
                            <td>3949107052</td>
                            <td>3949106658</td>
                            <td>394</td>
                        </tr>
                    <tr>
                        <td colspan="4"><b>Combined lag</b></td>
                        <td><b>394</b></td>
                    </tr>
                    </tbody>
                </table>
                </p>
            </div>
    </div>
</div>
</div>
<div class="pb-0 mt-5 mb-4">
    <div class="container">
    </div>
</div>
</body>
</html>`
)

func TestExtractData(t *testing.T) {

	var reader = strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	res := extractLagData(doc)
	if res == 394 {
		t.Logf("Success: %v\r\n", res)
	} else {
		t.Errorf("Expected 394 but got %v", res)
	}
}
