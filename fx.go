// 2020, GUILLEUS Hugues <ghugues@netc.fr>
// BSD 3-Clause License

package wsSaloon

// Function to call when a message is received. Take the input message and
// return a response (can be nil).
type CBReceived func(data []byte) (response []byte)

// Run the callback if it is defined. If the callBack return a response,
// the response will send to the saloon.
func (cb CBReceived) run(s *Saloon, data []byte) {
	if cb == nil {
		return
	}
	rep := cb(data)
	if len(rep) != 0 {
		s.Write(rep)
	}
}
