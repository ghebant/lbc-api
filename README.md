# lbc-api
Description

# Start the api
Start dev
```bash
make dev
```

# Tests
Run tests
```bash
make test
```
Run tests with coverage
```bash
make test-coverage
```

# Endpoints
| Method   | URL                      | Description                                   |
|----------|--------------------------|-----------------------------------------------|
| `GET`    | `/`                      | Check api's health                            |
| `GET`    | `/ads`                   | Retrieve all ads                              |
| `POST`   | `/ads`                   | Create a new ad                               |
| `GET`    | `/ads/{id}`              | Retrieve ad with and id                       |
| `UPDATE` | `/ads/{id}`              | Update ad                                     |
| `DELETE` | `/ads/{id}`              | Delete ad                                     |
| `GET`    | `/search?input={string}` | Search for a car model that matches the input |
## Ads
### Get all ads
`GET /ads`
```bash
curl -i -H 'Accept: application/json' http://localhost:8000/ads 
```
#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jul 2022 07:10:06 GMT
Content-Length: 270

[
   {
      "ad_id":1,
      "title":"ad title",
      "content":"content",
      "category":1,
      "automobile":{
         "id":1,
         "ad_id":1,
         "brand":"Audi",
         "model":"Rs4",
         "created_at":"2022-07-11T07:08:41Z",
         "updated_at":"2022-07-11T07:08:41Z"
      },
      "created_at":"2022-07-11T07:08:41Z",
      "updated_at":"2022-07-11T07:08:41Z"
   }
]
```

### Create an ad with automobile category
`POST /ads`
```bash
curl -i -X POST -H 'Accept: application/json' http://localhost:8000/ads -d '{
	"title": "ad title",
	"content": "content",
    "category": 1,
    "automobile": {
        "brand": "Audi",
        "model": "Rs4"
    }
}' 
```
#### Response
```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jul 2022 07:08:41 GMT
Content-Length: 268

{
   "ad_id":1,
   "title":"ad title",
   "content":"content",
   "category":1,
   "automobile":{
      "id":1,
      "ad_id":1,
      "brand":"Audi",
      "model":"Rs4",
      "created_at":"2022-07-11T07:08:41Z",
      "updated_at":"2022-07-11T07:08:41Z"
   },
   "created_at":"2022-07-11T07:08:41Z",
   "updated_at":"2022-07-11T07:08:41Z"
}
```

### Get ad with id
`GET /ads/{id}`
```bash
curl -i -H 'Accept: application/json' http://localhost:8000/ads/1
```
#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jul 2022 07:11:17 GMT
Content-Length: 268

{
   "ad_id":1,
   "title":"ad title",
   "content":"content",
   "category":1,
   "automobile":{
      "id":1,
      "ad_id":1,
      "brand":"Audi",
      "model":"Rs4",
      "created_at":"2022-07-11T07:08:41Z",
      "updated_at":"2022-07-11T07:08:41Z"
   },
   "created_at":"2022-07-11T07:08:41Z",
   "updated_at":"2022-07-11T07:08:41Z"
}
```

### Update an ad
`UPDATE /ads/{id}`
```bash
curl -i -X PUT -H 'Accept: application/json' http://localhost:8000/ads/1 -d '{
	"title": "ad BMW",
	"content": "super car !",
    "category": 1,
    "automobile": {
        "brand": "BMW",
        "model": "M3"
    }
}' 
```
#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jul 2022 07:13:31 GMT
Content-Length: 268

{
   "ad_id":1,
   "title":"ad BMW",
   "content":"super car !",
   "category":1,
   "automobile":{
      "id":1,
      "ad_id":1,
      "brand":"BMW",
      "model":"M3",
      "created_at":"2022-07-11T07:08:41Z",
      "updated_at":"2022-07-11T07:13:31Z"
   },
   "created_at":"2022-07-11T07:08:41Z",
   "updated_at":"2022-07-11T07:13:32Z"
}
```

### Delete an ad
`DELETE /ads/{id}`
```bash
curl -i -X DELETE -H 'Accept: application/json' http://localhost:8000/ads/1
```
#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jul 2022 07:14:37 GMT
Content-Length: 21

{"message":"deleted"}
```

## Search
### Search car model
`GET /search?input={string}`
```bash
curl -i -H 'Accept: application/json' http://localhost:8000/search?input=rs4%20a
```
#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 09 Jul 2022 10:15:01 GMT
Content-Length: 30

{
   "brand":"Audi",
   "model":"Rs4"
}
```

# How search works ?
For the matching part the Levenshtein distance algorithm has been used combined with multiple percentages calculated for every keyword matching a car model
## Levenshtein Algorithm
The Levenshtein distance is a string metric for measuring difference between two sequences. Informally, the Levenshtein distance between two words is the minimum number of single-character edits (i.e. insertions, deletions or substitutions) required to change one word into the other.
### Example
The Levenshtein distance between "kitten" and "sitting" is 3 because it needs no less than 3 edits to change one into the other.

1. **k**itten → **s**itten (substitution of "s" for "k"),
2. sitt**e**n → sitt**i**n (substitution of "i" for "e"),
3. sittin → sittin**g** (insertion of "g" at the end).

## The Matching algorithm
1. First **iterate** over **all car models**
2. **Normalize** both the **search input** and the **car model** (lowercase, trim spaces etc...)
3. Split the **search input** and the **car model** into **keywords**, then for each keyword:
   1. Compute the **Levenshtein distance** between the **search keyword** and the **whole car model**, then calculate a **matching percentage** from it with the **total car length**
   2. **Compute the Levenshtein distance** to every car model keyword
   3. Use that distance to calculate a **matching percentage** with the **car keyword length**
   4. Calculate a **percentage** of how many **characters in a row** both keywords have **in common** since index 0
   5. Calculate a **weight** based on the search keyword length (so that longer words have more weight when they will be multiplied by their matching percentage)
   6. Then calculate a **global matching percentage** of the given search input for the car model based on all the percentages calculated
4. Return the car model which as the **best average matching percentage** with the searched input

### Example
```
Search Input: ds 3 crossback
```
```
1.   car model :                                       s3
2.   search keywords:                           ds      3   crossback
3.1 .Matching percentage whole car model:       0%     50%    0%
3.2 .Levenshtein distance from car keyword :    2       1     8
3.3 .Matching percentage :                      0%     50%    0%
3.4 .Percentage characters in row matching :    0%     0%     0%
3.5 .Calculate weight :                        1.4     1     2.8
3.6 .Global keyword matching percentage:        0%   33.33%   0%
4.   Average matching percentage:                     11.11%

1.   car model :                                       ds3
2.   search keywords:                           ds      3   crossback
3.1 .Matching percentage whole car model:     66.67%  33.33%    0%
3.2 .Levenshtein distance from car keyword :    1       2       8
3.3 .Matching percentage :                    66.67%  33.33%    0%
3.4 .Percentage characters in row matching :  66.67%    0%      0%
3.5 .Calculate weight :                        1.4      1      2.8
3.6 .Global keyword matching percentage:      84.44%  22.22%    0%
4.   Average matching percentage:                     35.56%

6. Best match = ds3 : 35.56%
   ```
