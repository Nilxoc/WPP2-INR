# WPP2-INR

Basic information retrieval system which supports boolean queries.

## Running

`go run main.go [args]`

### Configuration

Format: `-key value`

The following configuration options are available:

| arg        | desc                                                                                       | default |
|------------|--------------------------------------------------------------------------------------------|---------|
| doc        | Path to the input file                                                                     |         |
| dict       | Path to the dictionary (if present, the dictionary will not be recreated)                  |         |
| k          | size of the k-gram index terms                                                             | 2       |
| r          | threshold of results before spell correction (only relevant if `correction` is `true`)     | 5       |
| j          | threshold for the jaccard-coefficient                                                      | 0.2     |
| correction | if spell correction is enabled                                                             | `true`  |

## Queries

* Simple terms
  * `blood`
* Junctions
  * `blood [JUNCTION] pressure`
* Sub-expressions
  * `(blood OR pressure) AND cardiovascular`
* Phrases
  * `"blood pressure"`
* Proximity
  * `diet \[k] health` (eg. `k = 10`)

### Junctions

* `AND`
* `AND NOT`
* `OR`
