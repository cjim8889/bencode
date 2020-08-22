# bencode
Bencode implementation in Golang

This is a bencode implementation in Golang. It is super simple and blazingly fast! 

Only two functions are of importance to you, namely func Parse(r *Reader) (interface{}, error) and func (r *Reader) DecodeStream() (interface{}, error)
