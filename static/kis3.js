(function () {
    var request = new XMLHttpRequest();
    var url = document.currentScript.baseURI
        + "view?url=" + window.decodeURI(document.documentURI);
    if (document.referrer && document.referrer.length > 0) {
        url += window.decodeURI(document.referrer);
    }
    request.onload = function () {
        if (request.status >= 200 && request.status < 300) {
            console.log('Success with tracking page view!');
        } else {
            console.log('The tracking request failed!');
        }
    };
    request.open("POST", url);
    request.send();
})();
