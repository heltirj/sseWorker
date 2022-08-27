function onLoaded() {
    var source = new EventSource("listen");
    source.onmessage = function (event) {
        console.log("OnMessage called:");
        console.dir(event);
        const eventList = document.getElementById("list");
        const newElement = document.createElement("li");

        newElement.textContent = `message: ${event.data}`;
        eventList.appendChild(newElement);
        console.log(eventList)
    }
}