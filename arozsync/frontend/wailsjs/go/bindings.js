// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
const go = {
  "main": {
    "App": {
      /**
       * Greet
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<string>}  - Go Type: string
       */
      "Greet": (arg1) => {
        return window.go.main.App.Greet(arg1);
      },
      /**
       * OpenLinkInLocalBrowser
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<void>} 
       */
      "OpenLinkInLocalBrowser": (arg1) => {
        return window.go.main.App.OpenLinkInLocalBrowser(arg1);
      },
      /**
       * ScanNearbyNodes
       * @returns {Promise<string>}  - Go Type: string
       */
      "ScanNearbyNodes": () => {
        return window.go.main.App.ScanNearbyNodes();
      },
      /**
       * TryConnect
       * @param {string} arg1 - Go Type: string
       * @param {string} arg2 - Go Type: string
       * @param {string} arg3 - Go Type: string
       * @returns {Promise<boolean>}  - Go Type: bool
       */
      "TryConnect": (arg1, arg2, arg3) => {
        return window.go.main.App.TryConnect(arg1, arg2, arg3);
      },
    },
  },

};
export default go;
