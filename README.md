# alexa skill
This is only a simple example for a alexa skill. The skill is written in go and can be run on AWS Lambda.

## run code on AWS Lambda
* You must create a aws lambda function in a compatible region like Irland for german skills. 
* To run your code on AWS Lambda, you need to create a zip file with all your code and dependencies.

## Alexa Skill Kit

### how to use skill
* start skill with invocation name
  > alexa open your-invocation-name
* ask question with intent keyword
  > sag wie ist das wetter
   
  > erzähle mir wie ist das wetter
  
  > nenne mir wie ist das wetter
  
  > sage mir wie ist das wetter
  
  > erkläre mir wie ist das wetter
  
  > frage wie weit ist es bis zum mond
* when you not ask a question, alexa will ask you for a question after 8 seconds
* stop intent by saying stop keyword to alexa
  > stop


### intent JSON
```json
{
  "interactionModel": {
    "languageModel": {
      "invocationName": "your skill invocation name",
      "intents": [
        {
          "name": "AMAZON.CancelIntent",
          "samples": []
        },
        {
          "name": "AMAZON.HelpIntent",
          "samples": []
        },
        {
          "name": "AMAZON.StopIntent",
          "samples": []
        },
        {
          "name": "AMAZON.NavigateHomeIntent",
          "samples": []
        },
        {
          "name": "QuestionIntent",
          "slots": [
            {
              "name": "question",
              "type": "AMAZON.SearchQuery"
            }
          ],
          "samples": [
            "sag {question}",
            "erzähle {question}",
            "nenne {question}",
            "sage {question}",
            "erkläre {question}",
            "frage {question}"
          ]
        }
      ],
      "types": []
    }
  }
}
```
