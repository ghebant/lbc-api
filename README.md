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
Content-Type: application/json
Date: Sat, 09 Jul 2022 10:00:56 GMT
Content-Length: 146

[
   {
      "ad_id":1,
      "title":"ad title",
      "content":"content",
      "category":"category",
      "created_at":"2022-07-09T10:00:34Z",
      "updated_at":"2022-07-09T10:00:34Z"
   }
]
```

### Create an ad
`POST /ads`
```bash
curl -i -X POST -H 'Accept: application/json' http://localhost:8000/ads -d '{"title": "ad title","category": "category","content": "content"}}' 
```
#### Response
```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Sat, 09 Jul 2022 10:06:41 GMT
Content-Length: 144

{
   "ad_id":2,
   "title":"ad title",
   "content":"content",
   "category":"category",
   "created_at":"2022-07-09T10:05:30Z",
   "updated_at":"2022-07-09T10:05:30Z"
}
```

### Get ad with id
`GET /ads/{id}`
```bash
curl -i -H 'Accept: application/json' http://localhost:8000/ads/2 
```
#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 09 Jul 2022 10:08:04 GMT
Content-Length: 144

{
   "ad_id":2,
   "title":"ad title",
   "content":"content",
   "category":"category",
   "created_at":"2022-07-09T10:05:30Z",
   "updated_at":"2022-07-09T10:05:30Z"
}
```

### Update an ad
`UPDATE /ads/{id}`
```bash
curl -i -X PUT -H 'Accept: application/json' http://localhost:8000/ads/2 -d '{"title": "ad title modified","category": "category","content": "content"}}' 
```
#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 09 Jul 2022 10:11:41 GMT
Content-Length: 153

{
   "ad_id":2,
   "title":"ad title modified",
   "content":"content",
   "category":"category",
   "created_at":"2022-07-09T10:05:30Z",
   "updated_at":"2022-07-09T10:11:42Z"
}
```

### Delete an ad
`DELETE /ads/{id}`
```bash
curl -i -X DELETE -H 'Accept: application/json' http://localhost:8000/ads/2
```
#### Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 09 Jul 2022 10:13:07 GMT
Content-Length: 21

{
   "message":"deleted"
}
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
   "model":"rs4"
}
```

# How it works ?
For the matching part the Levenshtein distance algorithm has been used
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
3. For each searched keyword **compute the Levenshtein distance** to every car model 
4. Use that distance to calculate a **matching percentage** with the car model length and the number of changes necessary to transform the keyword into the model
5. All the search keyword's matching percentage combines makes an **average matching percentage** of the whole search input
6. Return the car model which as the **best average matching percentage** with the searched input

### Example
```
Search Input: ds 3 crossback
```
```
1.   car model :                                   s3
2.   search keywords:                        ds   3   crossback
3.   Levenshtein distance from car model :    2   1       8
4.   Matching percentage:                     0% 50%      0%
5.   Average matching percentage:                  16.66%

1.   car model :                                   ds3
2.   search keywords:                        ds   3   crossback
3.   Levenshtein distance from car model :    1   2       8
4.   Matching percentage:                     67% 34%     0%
5.   Average matching percentage:                  33.66%

6. Best match = ds3 : 33.66%
   ```
