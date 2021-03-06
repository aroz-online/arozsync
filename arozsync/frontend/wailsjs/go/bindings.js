// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
const go = {
  "main": {
    "App": {
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
       * @param {Array<string>} arg1 - Go Type: []string
       * @param {string} arg2 - Go Type: string
       * @param {string} arg3 - Go Type: string
       * @param {boolean} arg4 - Go Type: bool
       * @returns {Promise<Array<string>>}  - Go Type: []string
       */
      "TryConnect": (arg1, arg2, arg3, arg4) => {
        return window.go.main.App.TryConnect(arg1, arg2, arg3, arg4);
      },
    },
  },

};
export default go;
