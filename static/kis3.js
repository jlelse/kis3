(function () {
    var request = new XMLHttpRequest();
    var url = document.currentScript.src.replace("kis3.js", "view?url=" + window.decodeURI(window.location.href));
    if (document.referrer && document.referrer.length > 0) {
        url += "&ref=" + window.decodeURI(document.referrer);
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
