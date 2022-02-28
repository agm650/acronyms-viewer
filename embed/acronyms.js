function trim(str) {
  return str.replace(/^\s\s*/, "").replace(/\s\s*$/, "");
}

// This functionn will verify if the acronym field is set
// If set it will request the backend for a list of definition
// Then it will display them
function sreachAcronym() {
  /*
  Fields must not be empty:
  - searchItem
  */
  var searchElement = trim(this.document.getElementById('searchItem').value);
  alert("Searching for: " + searchElement);

  if(searchElement.length == 0) {
    errorCount++;
    this.document.getElementById('errorSearch').innerHTML="Champs vide! Saisissez un acronyme a chercher";
    return;
  }

  // We have the data, lets look for definition
  // ajax(isReady, "get", "/search?name="+searchElement, true);
  fetch("/search?name="+searchElement)
  .then(function(response) {
    return response.json();
  })
  .then(function(jsonResponse) {
    // do something with jsonResponse
    ol = document.createElement("ol");
    jsonResponse.definitions.forEach(element => {
      var li = document.createElement("li");
      textDefinition = document.createElement("p");
      textDefinition.appendChild(document.createTextNode(element));
      li.appendChild(textDefinition);
      ol.appendChild(li);
    });
    // Clear content
    this.document.getElementById('main_content').innerHTML = "";
    defSection = document.createElement("p");
    defSection.appendChild(document.createTextNode("Definitions pour <b>" + jsonResponse.name + "</b>"));
    this.document.getElementById('main_content').appendChild(defSection);
    this.document.getElementById('main_content').appendChild(ol);
  });
}
