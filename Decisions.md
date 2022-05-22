# Following are the decisions that were taken when designing the system

* Orders were stored at different price level
* In order matching the traded price will be that of the older order (which was previously present in the book)
* The orders are stored in pro rata basis of order size

## book/parser.go
For the parsing of individual parts in the request string the considered methods were
  ```go
  newOrderConds := []RequestParserFunc{parsePrice, ...}
  for idx, c := range newOrderConds {
    if err := c(&o, s[idx]); err != nil {
      return Order{}, err
    }
  }
  ```
which bechmarking showed took *246.4 ns/op* versus
  ```go 
  // Loop Unrolling
  if err := parsePrice(&o, s[0]); err != nil {
    return Order{}, err
  }
  if err := ....(&o, s[0]); err != nil {
    return Order{}, err
  }
  ```
which bechmarking showed took *215.8 ns/op* versus
  ```go
	{
		// without functions
		price, err := strconv.ParseFloat(s[0], 32)
		if err != nil {
			return Order{}, fmt.Errorf("parsing price: error %w", err)
		}
		o.Price = float32(price)
	}
  ```
which bechmarking showed took *207.5 ns/op*.

The Slowest one was around 15% slower than the fastest, but for the sake of clarity I went with that one instead.

## Todo
* Build a reporting system ( report/report.go )
    * Event sourcing
    * Can show the amout held by them
    * Can show the trades in specified time frame ( some sort of timeseries db?? )
    * Calculate the total trading charges for them ( won't be a priority )

* Support Order Deletion
    * Was Supported in the previous version but removed due to added complexity with the PRO RATA system
