package parser

// import "encoding/json"

// type fakeNSRecord NSRecord

// func (ns NSRecord) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		Hash uint `json:"hash"`
// 		fakeNSRecord
// 	}{
// 		Hash:         ns.GetHash(),
// 		fakeNSRecord: fakeNSRecord(ns),
// 	})
// }

// type fakeARecord ARecord

// func (a ARecord) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		Hash uint `json:"hash"`
// 		fakeARecord
// 	}{
// 		Hash:        a.GetHash(),
// 		fakeARecord: fakeARecord(a),
// 	})
// }

// type fakeMXRecord MXRecord

// func (mx MXRecord) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		Hash uint `json:"hash"`
// 		fakeMXRecord
// 	}{
// 		Hash:         mx.GetHash(),
// 		fakeMXRecord: fakeMXRecord(mx),
// 	})
// }

// type fakeTXTRecord TXTRecord

// func (txt TXTRecord) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		Hash uint `json:"hash"`
// 		fakeTXTRecord
// 	}{
// 		Hash:          txt.GetHash(),
// 		fakeTXTRecord: fakeTXTRecord(txt),
// 	})
// }

// type fakeCNAMERecord CNAMERecord

// func (cname CNAMERecord) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		Hash uint `json:"hash"`
// 		fakeCNAMERecord
// 	}{
// 		Hash:            cname.GetHash(),
// 		fakeCNAMERecord: fakeCNAMERecord(cname),
// 	})
// }
