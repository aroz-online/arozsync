/*
    Connection Check
*/

function TestConnection(ipAddr, username, password, remember, succCallback, failCallback){
  try {
    window.go.main.App.TryConnect(ipAddr,username, password,remember)
      .then((result) => {
        //Return JSON stringify scan results
        succCallback(result)
      })
      .catch((err) => {
        console.error(err);
        failCallback("");
      });
  } catch (err) {
    console.error(err);
    failCallback("");
  }
}

/*
    Scanner Functions
*/

function ScanNearbyArozOS(callback){
  try {
    window.go.main.App.ScanNearbyNodes()
      .then((result) => {
        //Return JSON stringify scan results
        callback(result)
      })
      .catch((err) => {
        console.error(err);
        callback("");
      });
  } catch (err) {
    console.error(err);
    callback("");
  }
}

function OpenLinkInBrowser(link){
  try {
    window.go.main.App.OpenLinkInLocalBrowser(link)
      .then((result) => {

      })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);

  }
}



// Setup the greet function
window.greet = function () {
  // Get name
  let name = nameElement.value;

  // Check if the input is empty
  if (name === "") return;

  // Call App.Greet(name)
  try {
    window.go.main.App.Greet(name)
      .then((result) => {
        // Update result with data back from App.Greet()
        document.getElementById("result").innerText = result;
      })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);
  }
};
