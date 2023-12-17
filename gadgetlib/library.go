package gadgetlib

type CommitMIMC struct {
	CommitString string `json:"commit_string"`
	Witness      string `json:"witness"`
	Message      string `json:"message"`
}

type CommitSHA256 struct {
	CommitString string `json:"commit_string"`
	Witness      string `json:"witness"`
	Message      string `json:"message"`
}

type CommitMerkleMIMC struct {
	CommitString string `json:"commit_string"`
	Witness      string `json:"witness"`
	Message      string `json:"message"`
}

type SigECDSA struct {
	Signature string `json:"signature"`
	Message   string `json:"message"`
	PublicKey string `json:"public_key"`
}

type EncrypAES128 struct {
	Plaintext  string `json:"plaintext"`
	Key        string `json:"key"`
	Ciphertext string `json:"ciphertext"`
}

type EncrypOTP struct {
	Plaintext  string `json:"plaintext"`
	Pad        string `json:"pad"`
	Ciphertext string `json:"ciphertext"`
}

type EncrypAES128GCM struct {
	Plaintext  string `json:"plaintext"`
	Key        string `json:"key"`
	Ciphertext string `json:"ciphertext"`
	Iv         string `json:"iv"`
	ChunkIndex string `json:"chunk_index"`
}

type SCsessionKey_TLS13AES128GCMSHA256 struct {
	DHSin                  string `json:"dhs_in"`
	IntermediateHashHSopad string `json:"intermediatehashhs_opad"`
	MSin                   string `json:"ms_in"`
	XATSin                 string `json:"xats_in"`
	TkXAPPin               string `json:"tkxapp_in"`
	IvCounter              string `json:"iv_counter"`
	Zeros                  string `json:"zeros"`
	ECB0                   string `json:"ecb0"`
	ECBK                   string `json:"ecbk"`
	Key                    string `json:"key"`
}

type SCopenRecord_TLS13AES128GCMSHA256 struct {
	Key            string `json:"key"`
	Plaintext      string `json:"plaintext"`
	Iv             string `json:"iv"`
	Cipherchunks   string `json:"cipherchunks"`
	ChunkIndex     string `json:"chunk_index"`
	Substring      string `json:"substring"`
	SubstringStart string `json:"substring_start"`
	SubstringEnd   string `json:"substring_end"`
	ValueStart     string `json:"value_start"`
	ValueEnd       string `json:"value_end"`
	Value          string `json:"value"`
}

type SCopenSession_TLS13AES128GCMSHA256 struct {
	DHSin                  string `json:"dhs_in"`
	Plaintext              string `json:"plaintext"`
	IntermediateHashHSopad string `json:"intermediatehashhs_opad"`
	MSin                   string `json:"ms_in"`
	XATSin                 string `json:"xats_in"`
	TkXAPPin               string `json:"tkxapp_in"`
	IvCounter              string `json:"iv_counter"`
	Zeros                  string `json:"zeros"`
	ECB0                   string `json:"ecb0"`
	ECBK                   string `json:"ecbk"`
	Iv                     string `json:"iv"`
	Cipherchunks           string `json:"cipherchunks"`
	ChunkIndex             string `json:"chunk_index"`
	Substring              string `json:"substring"`
	SubstringStart         string `json:"substring_start"`
	SubstringEnd           string `json:"substring_end"`
	ValueStart             string `json:"value_start"`
	ValueEnd               string `json:"value_end"`
	Value                  string `json:"value"`
}
