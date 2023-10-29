import { jsonConfig } from "./json-config"
import { __jolt } from "jolt-demo"

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

        return JSON.parse(await __jolt(response, jsonConfig.spec, false));
    }
    
}

export const hltb = new Hltb()