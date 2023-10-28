
export class JSONConfig{
    JSONRequest(param: string){
        const params = this.extractParams(param);

        return JSON.stringify({
            "searchType": "games",
            "searchTerms": params,
            "searchPage": 1,
            "size": 20,
            "searchOptions": {
              "games": {
                "userId": 0,
                "platform": "",
                "sortCategory": "popular",
                "rangeCategory": "main",
                "rangeTime": {
                  "min": 0,
                  "max": 0
                },
                "gameplay": {
                  "perspective": "",
                  "flow": "",
                  "genre": ""
                },
                "modifier": ""
              },
              "users": {
                "sortCategory": "postcount"
              },
              "filter": "",
              "sort": 0,
              "randomizer": 0
            }
          });
    
    }
    
    private extractParams(param: string){
        return param.split("%20");
    }
}

export const jsonConfig = new JSONConfig();