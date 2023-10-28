import { jsonConfig } from "./json-config"

export class Hltb{

    private url
    private searchURI


    constructor(){
        this.url = "https://howlongtobeat.com"
        this.searchURI = "https://howlongtobeat.com/api/search"
    }

    async search(param: string){
   
        const response = await fetch(this.searchURI, {
            method: "POST",
            body: jsonConfig.JSONRequest(param),
            headers: {'Content-Type': 'application/json', 'Referer': this.url },
        }).then(response => response.json()). then(json => {
            return JSON.stringify(json, null, 4)});

        return response
    }
    
}

export const hltb = new Hltb()