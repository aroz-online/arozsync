import * as models from './models';

export interface go {
  "main": {
    "App": {
		Greet(arg1:string):Promise<string>
		OpenLinkInLocalBrowser(arg1:string):Promise<void>
		ScanNearbyNodes():Promise<string>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
