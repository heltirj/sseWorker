function onLoaded() {
    var source = new EventSource("http://localhost:4000/listen");
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