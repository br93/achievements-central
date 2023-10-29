
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

    spec = JSON.stringify(
      [
          {
              "operation": "shift",
              "spec": {
                  "data": "&"
              }
          },
          {
              "operation": "shift",
              "spec": {
                  "data": {
                      "*": {
                          "game_name": "&2[&1].&",
                          "game_image": "&2[&1].&",
                          "comp_main": "&2[&1].&",
                          "comp_plus": "&2[&1].&",
                          "comp_100": "&2[&1].&",
                          "comp_all": "&2[&1].&",
                          "review_score": "&2[&1].&",
                          "profile_dev": "&2[&1].&",
                          "profile_platform": "&2[&1].&",
                          "release_world": "&2[&1].&"
                      }
                  }
              }
          },
          {
              "operation": "modify-overwrite-beta",
              "spec": {
                  "data": {
                      "*": {
                          "comp_main": "=divide(@(0),3600)",
                          "comp_plus": "=divide(@(0),3600)",
                          "comp_100": "=divide(@(0),3600)",
                          "comp_all": "=divide(@(0),3600)"
                      }
                  }
              }
          },
          {
              "operation": "modify-overwrite-beta",
              "spec": {
                  "data": {
                      "*": {
                          "comp_main": "=divideAndRound(1,@(0),1)",
                          "comp_plus": "=divideAndRound(1,@(0),1)",
                          "comp_100": "=divideAndRound(1,@(0),1)",
                          "comp_all": "=divideAndRound(1,@(0),1)"
                      }
                  }
              }
          }
      ]
  )
}

export const jsonConfig = new JSONConfig();