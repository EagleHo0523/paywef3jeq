package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	util "../util"
	v "../verification"
)

type respSetup struct {
	Public_key string `json:"public_key,omitempty"`
	Password   string `json:"password,omitempty"`
	Token      string `json:"token,omitempty"`
}
type reqVerify struct {
	Uid        string `json:"uid,omitempty"`
	Sys_pubkey string `json:"sys_pubkey,omitempty"`
	Token      string `json:"token,omitempty"`
}
type respVerify struct {
	Code  string `json:"code,omitempty"`
	Token string `json:"token,omitempty"`
}
type reqLinkWallet struct {
	Token     string `json:"token,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
type reqHistoryWallet struct {
	Cid        int    `json:"coin_idx,omitempty"`
	Start_time int64  `json:"start_time,omitempty"`
	End_time   int64  `json:"end_time,omitempty"`
	Token      string `json:"token"`
	Timestamp  int64  `json:"timestamp"`
}
type reqTradeCreate struct {
	To        string `json:"to"`
	Cid       int    `json:"coin_idx"`
	Value     string `json:"value"`
	Expire    int64  `json:"expire,omitempty"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqTradeCheck struct {
	To        string `json:"to"`
	Info      string `json:"info"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqTradeConfirm struct {
	Tid       string `json:"tid"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqTransferConfirm struct {
	To        string `json:"to"`
	Cid       int    `json:"coin_idx"`
	Value     string `json:"value"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqOfflineCreate struct {
	To        string `json:"to"`
	Cid       int    `json:"coin_idx"`
	Value     string `json:"value"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqOfflineCheck struct {
	From      string `json:"from"`
	Info      string `json:"info"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqNolimitCreate struct {
	Cid       int    `json:"coin_idx"`
	Value     string `json:"value"`
	Expire    int64  `json:"expire,omitempty"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqNolimitConfirm struct {
	From      string `json:"from"`
	Info      string `json:"info"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}

// type reqTradeTid struct {
// 	Tid string `json:"tid"`
// }

func main() {
	// toSetup()
	// toVerify()
	// decryptToken()
	toLinkWallet()
	// historyWallet()
	// base64encode()
	// decryptoResponse()
	// createTradeData()
	// checkTradeData()
	// confirmTradeData()
	// confirmTransferData()
	// checkOfflineData()
	// confirmOfflineData()
	// createNolimitData()
	// confirmNolimitData()
	// toLogin()
	// backLogin()
	// fmt.Println("    Unix:", time.Now().UTC().Unix())
	// fmt.Println("UnixNano:", time.Now().UTC().UnixNano()/1000000)
	// calcWid()
}

func toSetup() {
	tmp_pubkey := "04bacc176de5a56d64021c698141714ffd3911eddc995209f2d0aaba599ac74d579e4ac30c4738eb865f0ad1863c20200dad72a180d2d0610212b3eab9fb28ea67"
	token := "b9fef2cee52e8f83cb75bbe93094fb2b8090ef93"

	psn_pubkey := "047a251a43d724819147c2bfeabb0e04dc6e37d4f30a36ded23c7d91254e50dfdacc62f9060ba5107c6f81d1a14b9c7bcc49f8fb64c6f8ebfb6b7544db4ff2bc96"

	s, err := jsonSetup(psn_pubkey, "password", token)
	if err != nil {
		fmt.Println("err:", err)
	}
	veri, err := v.ImportPubKey(tmp_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(s)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}
func jsonSetup(public_key, password, token string) (string, error) {
	// var rtn string = ""
	resp := respSetup{
		Public_key: public_key,
		Password:   password,
		Token:      token,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func toVerify() {
	s := "c54d9c7307812efd5a1632e599d79ddc02ca0020f4c2b90a52eec0e7ad28231ab7792db2dbede01c15212c2d02c36ec155697b2b00205872ea7c63605a4f7dcbff2094297d0f3adc679e8f42cff8cb155f0d5e278c3d84a77df39d75bdd6ce7e5d05887534be3a875de71d0ff4357024cd36a86590302bb668f2f4e6f11dddb2bda7a320f7cc8b0482893226ed6f5e2e8b27573807e628f994534c5856fedb014bc9b72efc3efd39571ef56b8d645317cf05fb7ae986d0158157af22cbe37bbb0c5de08f9fcd0c4770dc8066b0694c3740a0186dd18acdb8bc0fcf972f504581a86125fcaa9a0a7c0c2bea3c22615afdf0235edfa656657fb45a616904d9d175c329c34e80c1e185555a67f8b287b44657e15f3b596684e0fa209bcd5f1db22e7fb737a4db065d0b73b90e9f9e6816bd2b2f60ab98b1aa79648b22f5fbd33120502ff3727e10edfb789c145fab7a9f81ef43e3a458e4af013167d69eb735d972168811a723c1a9ed491a0ee2c23b52d4c317f059954f"
	psn_privkey := "5ceec8d7f5e0cbbb082431408cdbe131068b9bedc87253fd5d30db1424b6b47a"
	account := "0x1F6998A6153378aeC050Ca02ec37B2cFC6ddDbF7"
	psn_pubkey := "047a251a43d724819147c2bfeabb0e04dc6e37d4f30a36ded23c7d91254e50dfdacc62f9060ba5107c6f81d1a14b9c7bcc49f8fb64c6f8ebfb6b7544db4ff2bc96"

	veri, err := v.ImportPrivKey(psn_privkey)
	if err != nil {
		fmt.Println("err:", err)
	}

	ct, err := veri.Decrypt(s)
	if err != nil {
		fmt.Println("err:", err)
	}

	var req reqVerify
	err = json.Unmarshal([]byte(ct), &req)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("sys_pubkey:", req.Sys_pubkey)

	st, err := jsonVerify(account, psn_pubkey, req.Token)
	if err != nil {
		fmt.Println("err:", err)
	}
	veri, err = v.ImportPubKey(req.Sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(st)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("uid:", req.Uid)
	fmt.Println("encode:", pt)
}
func jsonVerify(account, public_key, token string) (string, error) {
	code := processHashCreate(account, public_key)

	resp := respVerify{
		Code:  code,
		Token: token,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func processHashCreate(text1, text2 string) string {
	s := sha1.Sum([]byte(text1 + "&" + text2))
	return hex.EncodeToString(s[:])
}

func decryptToken() {
	r := "1003b54ab59a2ce2ead6bc0cdc6671f002ca0020952c6317de1730c8af0ec09ea5269169344815c75b1e48c96c821314b006bfc10020f50f97a6297e2deca94b1f67fa9a0535dd57962a426ab8af7a5d31ad7de4048f102cbca52de7993c1777f8af210ac9faeff4d2bcec2abdefd63c220db0574ce765bbd1eec1f031170115236e8bf26316b569dcce412dbb641482b269d1628aaefd9814334040b9735eeab4dc7a8dc7da50eba77cf33db857805e926c09fedcf00f25110ae0b87e51f1bdfb9ca4857ba6"

	psn_privkey := "5ceec8d7f5e0cbbb082431408cdbe131068b9bedc87253fd5d30db1424b6b47a"

	veri, err := v.ImportPrivKey(psn_privkey)
	if err != nil {
		fmt.Println("err:", err)
	}

	ct, err := veri.Decrypt(r)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("ct:", ct)
}

func toLinkWallet() {
	sys_pubkey := "04a61464e5d03728e46848e2f2472f32c004c2f5e0745f20a4537b640b651d3b4f40ead137333c718dbaa1054058ad2c0dae050f603e49c1ccbb318aebffc0e5e5"

	var req reqLinkWallet
	req.Token = "c3024b6cff917ba4a85bd41f389057b84be72734"
	req.Timestamp = time.Now().UTC().Unix()

	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}
func historyWallet() {
	sys_pubkey := "04a017ac0e9665ed7dd02df71ace98fde3b2a6a109f4181c0e397392a5e353d36d5b6d7301fe24516058ee5bca4172dfa0bdfa78b98fd24eb4a58da0441f665731"

	var req reqHistoryWallet
	req.Token = "c3322ab5d9bae0daaec69ada70fa2547a6ff587e"
	req.Timestamp = time.Now().UTC().Unix()

	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}

func decryptoResponse() {
	r := "af9b4548c123283905ad4d830cebf8d102ca00207e902f5d9e0bfec4505d2b921ac6c832afddd9d8c9522d7463062a425979a094002038cbc12d19c2101f01fe15072a428369f914e09de3259b9c5271704b078384cd285f294a456b528a4763371a51158cbfbc97d54f8d3f7164b5bc39f1f3c8cd8af4667325f5491d1e927d1d02f18c017125ae3ecc1488fa360b709d261306e2d7fca79d6eedf7e711965cee1333b1aacaa2f5d6833246e323e5d7ddf4e2e29412cb673f2693d730f05b3ac4ec5102bc6bfc47c6e1c38e1b493ac24bfc7a52119d9c0d0245c692a7a543afc64b1cbbb0a8d1aa5b3f91fff91c9b4b81de2b63875cfa916ab29df3a7b27d5b28155b7890e03b60f66ab4523a60e58efefca15db7878bf55b8fc34efb306f9fcc54ecd25fb42b3165934b20b1f10b59b351df647b3632790e3e0b94b28fc0edda2a8c2e891f62ef3b72b5bf3bcb295433cc1d2615160b0c0a5fc8ce216a2a65549862dd1c24a98f2d823612a73cab48ceb11d8a15bf9f77fdf5b113db5aa6ce59aef714cc755b0a08e3cab6cb38a161b5f8d864efcf87c5fb407fefe0b5d7b302c7adac8d22ac16a21d45a552a866e4368b708a90e8127e6068ab63ec2e28dab9af634a9059c91ff196df62dbed4af3d30f0a6c7e4c163815aa558081c5e7ca114999792a5bc09cb70e804d1905008267e38aef34e6ab8059c91466f2aa002b595bd3a12079852c792e40ef6e0437d2e452afefb4e35e883c441bf1b73fa20579ad4c84a15101dcf78f6af552d40120ffb3e83e9c2a9954f67d3a15a763dae0e7533a02b79d6a31a021f39c74187912e662d13ef73a25b3f1e107e68b1de2cf026bc8ee189b7adca740617dd10939e3b5742225800798303dc34fcac1647436fae8fd6d219b9f667d08c592651442bda3b2cb02b186c78628df1ba0acaa0c185fe0b17431522d41d5709b5754e124eaa9bdf15baf9daf3adb68a684ccc8330bf65786935b892de749cfcf387f38e1e7318e593ea7b56ccc6d72c3e1c4f7d28024dc711bc5fc86b4da6d2b7ee7e33d81cea22dac701dab05410ba2a2a984791f20e7bc4875baf3f4995df2a190fd8ed2d929e17cccf5b7e3e474abef86bbb112763fffc977364c2cc4bc38e551f7e4e91e17bc70c937df0ab3913415f78f2eac907449edb3716727d3902fb895f489111eb7d2a306a14cb474ee714260b867f3e3d5474014973f809154e866cac58b0a19bd6713e82c38c1ea880ac97166dab39ae611f3a3f15c3cc21eb157859fd8f5bf7d70f2e38f9fb723804d3601736961d718ab2c1f20de7eb7e5e425dde9e347cec0f8761f6ecfb39aa560def3c89575aab29581a9f9bcbc6116fff36889d65697ba88ce8c7b4cfccae9fb1c7860188bd5f4cf088b8252da05793f3cbf0705e67879aec1da4b881f0674db39e958ccaeafb91df2dc8d488f6666e73499731b47d45a7c788cfe29dc2e5681bf40c1be354c9b3ab15c2d925d896a409bcf4dfd457880bc7e50e5bdddd54a273b31cf5107e945da80d66daab1cf4702a620cb0d8af7e90033233d657288f6c741acd1919cb27b0e0c2149c72bd44376f75300fdedef78fab29e32bcb4a85813179ec8512a21a484b2add1e2d8b556ad21a5a9a23237d1cb77ab8e94aff20c9c4045e3a1524841a5fe3dd636d79c3d51e9f8e79208167f104193fedf00fab76984ffa081d061085280ad0fa6c1dfc031e075527e4e80d1571d799467bee4439f29bf554a42445522a15d5fb1cb8dffd58193c1b436dcc37c65ebd0939aa7950239be36de6f257dc3d53c82d37a77b4ce38a44af23d628ab4004347f6da9ad85a88c1aeff2447540729acf42ad89b3122d5ef6f9627f468aee3b6303e74a86ce0a74ff9e959602164eee7e515bd464ac5540b28715152f46c8d05f73f9e35f207a8c5702192ab3d2d04d8ce59d0cea01e8ff45aed45061ead53e16558214aa931f68eec98d09b920f39f747a8c2328c77b753e503cd629c5898bf67c451c41f63e274f1964b13f1739e1540d88b80f46a091a4010797e125daacc5a7e2dbc9495d5f91a067c790d93d3b224269a01df39f45029e7c9f56d0ed76fb4c40107bac361cd3a6b03d8649c2bed47d73d5c429c8d380772b88a059008b5bb1e83adb6f80ceef53eac603fde8ca139e965d0fbce37b246f43fe2d3b539d6bdc956d8d2df00bd0a1206229dcd271d06bf63e416842698ec997c1d3eb1ed2d2449b48ea84cb7fdfe36f908dd9c1a201897400e643c0092f605ba81ef604ef8975ed137ac0a72a326b42612c517e418ff75ee939c383e61006fe36f1a3db4123c7b78df23f9c14daccfb138773c8166a00c94dc4778f4eb552bd7d6c10e512ed625ae3831feee464b96062632615d43891fd3d07f3b37bce990281cfa2faa2cf0fa7dbef74a2cac5a2ef0d4d4ed2449a964dcea73471f8f0b0cc8dc13abbd451e7523297f23bb508597c5508a78ebf291841946657bec9f4e4d1879f52666c8fcf0e31f28e3fecac9682b0d78dbb4af10d1aa4fb1429245aeb2eb7fbd484abbe442e34d0e82d0ce078a9d33bc87687f0ce87ed9e15c0c0e0cb9b7f2c5838700032a9c04205cec0fb8d567e89edfbad0bd31cb97c07e3b6d9f9bd413f9fd90baa49f052a4682ecc1c466c736ce42768245aee01c6d7d059f947cde921d80daf68ba895f06f3cdc691501343da2322574644beb204eb6b6e913e900ac0ff408c5b4fa7b46576d9edf364db068c9b63c0233f5f993d087c72a8fd786213bda21a2dd7e7467fd245f3910f30cf9f329d8edb9ca2ad5d0ecdd518ff80039c1bcfa9088f52bec5ba37c76264476d21dc2731e5c7d14255cd1450fb97f372ec464fe693d14f70277fb9bdf00c2380417634d0da826a710a15cecbbcfa4f0af48825f68c6946ac0a1d9643eb14cf7db974ac787267ea7705e322ecae00e41fa0886e56ca93344ead5864740d78e0990d7b91aaabc97735936ba22d01fdfe8f368f502ca39591c48a4666c5bbfbfd6b8afff300deca18936778bdc9dff70aacef9ac4ab36da5a069ecb8f37e86befb7a590d3ad34ed3ae549e508fe0850bebf77e913fce452bb2ea2fc76b4297b7609a339397ac8bc73c645d87feb2312e90fb053f545fc490a74fe55e3d7d4bdde57bdea7ff0a327fe14d4056dcda0a513eb4e6b8e80abec7d249fbd0fd1d1b5713a5509a7d9b30a55b002917376410d3d04008e994e807735a61324a908216b7fcd02852e93b3e7cb66a43c201ec8e64144958d8ac192d45de6771583fa88b4b066cb4b78040777c9524972bb0271f7b5f86dcc6350ec7691a855a8d42834bffaccf3241963cc98f3d1fad2a2765d714ce265ff3cdbecc9920acd2ed8d790a3e7a1d55278f1e14e986ca8c29620a2060c7bdb6e45ea8268b58a71725d71d1cc77536fc2e6fbd8466bf660fc27ec20f76baa2769bfd041c3d48d3430466f93a4dcc29a2c05f8d3a8136c1d83cc896b8b32555f9aebecd728126f25faa351548e86e44d51881e7abc5d194ed1ddfd1b48c383aca75936a8b387e5fcb6a96ca76cc0030aa3a797fc66fbde68e889ea6427563fee05e4881378294a9d369db902571f762b0007f9454cdc558c7ad8d8088f5eb7d61034b503f613b6a492950756f7b4cf9482949b008dc955834b5bce366052131b6531701efbf8baf39af5e50c2ce35c6c6a2493b668b8c271c2bd716392875fd2c037a6e5569f3a5c32e9b79235ba33f1c232f6cb9e512a92e5f7cec9bd6dac9780820a737547ee12357b730514f351a80aba8aee7f23c5a892bb04f05cd3af66e16c991c9e86575453d838fcbf65f5884d984e6d8fc79f3ac7354059a717a92091ee808fed0fe904e1fcfc0c3a5ed16916f537ad255fb2045b34864ac1e67e998c379184b7ec02cfb7b0dd58cd0373e6ad55c8f3f84eb4d9199dc64561a2184d980765463f35ea6d21dc850843172aab671ab9afd852aead3bf45df59bc7368ef3398d4e6a707104fc0a8b3988efd5ead3211c96e6cc3661fdff02bce7734c05c3edfd6d7dd64acfd944058cbb01cfd3e90a9577fc9ecc6c787f673320d1aff18759d49e3da196e3af2e00d9a7c386ab6119b043efa37df6c1e247538daf6e59af0b9f4327386836a007ccc11a5789ee46d9d478b663fc7ef52ec977d01f88e1b7d22d3bc7700d602eaea26bf60b158bcca11f7c92e64fa9724fa4e2b308d1e2b028dcd92ff2aa41d2f2401f26018516c4d0bf5d703301a957069f4843b29c71eaa47e159f66540fbbfd9a6d82b2908371baa79ba7d380a51f2f16f501ce876a958587e1cbae95f4ae6a8366dbd9c2c23f27a4b07b4a1ea03c0cd82acc0482b21fd3a6f8ff86d1aa1faa39e1eb0e290ac32d41ab5fb1f0b16cd8c3433831ca5841b5364daa1ba1587ac2264f896e5eedabd328d7b639104335306f298e2906402354b5a1fcce105af88987f54ec113112d0149da1dd05177fa2b90e96a72bb9c47844bc02ba8f1c0f6cad2b360016419fe99689f0ade623909709ab8ac4f952b6023b1780c36546abdbca9412760d5a58210f2c6634194eae2b08f26e38b2ab7475274cf9f29d509d0bb480b708d8d0f615c2ff94af30600e5917941c2e9e72488ee39547c184e250923d88189a1dae3faf8630c26850fc33e1316d84f7612b9bfe1a60fbbf29e9db2b515325650818aff95201c5bb623077a0b8718660418677a597573332ec833e411ab4ace53b5983675f57cc1de0a58c438ddc4dbbd46de760e353020276895a10af5497ef6c00f5b7c922f900bb135386d1e0322fb7e3cf91d4855656e306b33ae198d320e8fe85233a6ec5958d4d7e06f731be3f3b05da11d251662d9470a68771f97c0081b654ceb2bae56123d5167e2d7a30e60f8e577e8ea3d6d8e07bb435c8d71be6a33ba828edea285d60839750425ef8bc7f784e32b05461c11f177721fd7561a4d4d1b2f44781b6aea6c3a71a67327446b5b9539a7e192b84b34a8fd659194e28a659257a9bd7238c40c802ed1af6be69ef595fa17b3e8f1d8c3d5e39fb97c8a7728104f58b7c17594b8a01b10e76864cd55153d90c82bb9e1b9ea47eb0480785f180a1fd88d6941b0fd8661d94c78f05d8b1b1a327e46a7a31f28ed42c1327019b9d8adb67957c1a7d67c7e66905a99df5f930f4180df673f498dc0a20aa19f95b56b8df0400b854f48854ce6d0ed3625ea6577ac096df155ce5c5e9a09b696df290b16ecdb09b78424384bca7fdf5ca2576881f541fed1abdd7ca1134c02125749103b00c4c7e3cbc3a7ebc01bd0c92ea23a75dc2d898eaea63fc6b2b80ba41418a51313ac0c1d5c00d404e22a82cb3d4eac6e4469281b87cc5af3b840ea04f103bb7a8d9fe271b56b584052d9e901acc89d97e5ec3a58a710213ea6ad8f039c35492259c450eaadff879d47d59c796c621b58b20d38235362c1278cda4c3c2cf47782af24b626817ac1a5422d45197c7f418acd479f7239b555e232ca1c3b07641edf9c9e956cb34f4e672398465ffc85fdc2560367726e7f6f7b96c34895c87a32623e0921d1b1b20a04dc80b2c3508f87fc9bd837951859fd182eb68ac30394563c1186e5c9dba27ab76b96a2af19b126dcc880f1c05da3c9b9e6450384cb8ec81c33162d91c8c596f162f90b7a4077052fd5ee172de30ac4f2324c918390ee287feaa221c4989490eab5a14eeb46cbbca8034a51705150dc39565678a64c9fd2a3b1c14e69da7eac735875589163ea5fa7a4d9b7873013589aeecad84f4c7d4e8a6679e9b6c0cb5a0a0665edcc74264276a6f80d5b5b268b72b6a83f6cfc2b4e0997b775c55302ba5af09314b0c514e499e7393c226c69ce55fd56b607f17be0f23b8ac57a9c3127b04b1dac17b3acade5edeb558eeb1ecd97ea5a125c13c22d57d49e50b149e523296c279c22f6e219e2d07c0487d496e2543de858725405fecf2afd00775d623e9baae9dca35c255bcf2fafed4571315e5876c5ff4655ed4e9a277b2bf7692709a9e1cac1ad68ae9a00498c4be728ff9de15a5986b959557259db6a8d77b4882d7cad594f911b26bc983e045033ec38e0755cffc07b7cc8a3bdb1096d79a7dc247f1755b6b606fb45e0df65956c94b0b7bb3d79ac0c6a6d18441f29d5402d864c6f2102454f202d948beadac3503ae671fc9ad0a02e62011119951731a3d697ba4df87655aa5b2296f50bccfff3731fab4ab782b785f1ffbed8150741ed471907c21e2e3990975ec629cfe9de2421ef7407192ef6f33ffbd4b679b94cefc13e234066e77abf269a21824d69d9fbf69f1d71ec2bf48382291dede154b0103ca25275c050ddceba274a8d819ef218ca91cc77355cb1bbff38262128c6cf70d8fbec5150b2fcf2cb675a628b23769badeaf64740fe0b8963437fa66f8f64dc497e35797f112a6784af9cb8c287b67149392f44b829235bdc35d092b3ec0750299111943dd01f9892c80d2369c0d782aa03fc0f11e42055db22b111412d031fd8fb9a443d7102ec44455531a7d0735a74231168a8f8574017dddaab0ce3cf1e197fc408ce8ae6c022f69ba7a915712a6d5925b0ebfc0c962dfa5eed8583d6c807cb79530e32940b46a8cfd4381142344eb566d21fae8cecd1b3774dec91577f47374af7e92b0bdbabe96439cd23f22cd3d0f5dc241cf94efaae68af2cae34c57aa85013f4a2897082957cd7fca0446a1832b3434fcf34116c0449bd11e39d5d2b2ebb01b7a836acda1ac7d07ecb4a0b353e155046952da26fd7b3ef61a5db98bd8b1985b512391fdf0404c54c7f50decd558b3e3d804099ea90d4a0ad5b3811ebaa96881774a5c933f8ac5855ffab6499ea1b96a217139cf8f48fcc00cd2c131d6871dfda9547ad75e6cb02cd7be07afaabf1d03bc95963d29e762650f03ce39192d58d85efe10744ff5521764eae7ed1efb930f1dcbefac8d2ab1c3891f4118415fb6f4b48512af279dc1f4d9b1f10de26ddb4efdf852f711370f1262dd8ebd58f115854a31ff239906cc7bf4d6e153b7a503dd45b6565d7100bf31cfbf3863ffe9a8e1ce28b1e4eac486565c9385e7c6faa9dce25bd39a3785be06ad9473e30c1a68c2b92a21d2c0ac25f9cc694c7c285d7a2db4a066b3efa3a2df12d641ef3862851d38f5f8c54da9990ee717eaa86daeac245be72fe163a3767c69bd208b8b19ccaad8dd906f4949b1d406bbf574610512b30dad411d6de8530543889cb4e1acfc8de08d038f0f214337328b9f789262027ab0bcf142c3889330d8060f187e114033bacda65e3bf4c425ae21dd362fb0c82e7997ec0b392a15154bff097f0128ae0f08fe6ca818013aa808c16d3d45041a21eeee3c8aeba5b7d0c121e095e257daed3f705a0ca3e442d0d473614d4fc70b48198b0a030e889e5f2b78751693e507b062d3dd5964702fc96d3ef6f00e6f6b4a2fd43779c600f72d387cee96072d89909485840969669fad1f33f4e137bd86c6bfd643fdbcc758665589f72e0fa09a0dc81a4516449565012cf8ce3f6cb12128c594a9bcfd603c8e83e145cd19ecf588b1b73bc651ded38e6013e127bcdaf63ec91c608e783aff46656f0f6affff2b48b781a254de347ddfbeb98b67be22d11b3da2c4c8869078f1542b2aa8779024012c8ea74207d2cbbfbdb40a6d5266f8586c60b46827c26d8a4981f2fb1778dbe90923715dae077ec28c060e619e8ef061f1a397df0941d0ff739d9bff3837ad264baa62f5a8d0db8305175f0f19c17514d23b2ca23cd2f36091a6c23c516b302352592dc60709da72ff4ed667a9f424bd7c6dee996644d9f283193aee57af496bceb5d4e874531b4fc78de9a4b5f3e6b759eda2aacb5a802ffa97f401df01367cff3f86448a495fa523056f67b7d1b6e32356cc6170d38773b6c09785364bc77c27325cf7820a7a4b027e8897b8695a4fd7609e51fbc681924be48dc00e98b3784dcc1949cd7ad64e5a2028e166126233ef31acdbe469f24b1ad30d431ee428ef74882130f54759c88280701061f369f690877d2e856342e80213059beb5e43f7a761c4f73ea0f6de88a4815c244768a447c0ec9300615e69a80062f969f74cb94952c3f5a081afba342646b39b6f31b4bf30b92d09fe940f9fa93c548ff6ce7e4a85eb8f7c38fde17e884e3c2806ab2abcdf9b36577da97a63698db8e0d19475a61d90df0518ef1eabbd48aec08bad4960c14fc59be14c0d5fc58356d1478c4f44f62d752ced1e8a0bf34edceda27b223a14ac3df5532786e0338e7ad8e3e229d05b25e205b2b55178684d96ca0bec3a24abbb4bca52ba8710e055791d87fccae5dc630e2b97407a533b4ad8b295166922de228057af46396f12263f2f89231a18645f96be998dbe82b6a57e5c9f7475505af5dad52414a598a3515cc0543b836278cb348bf52cea7af9ee1104892730712ac5574907c7980050090d71721fb473c6a0039a26f97d42a07bc41b4f77b04387826ca5d1f04ea583557eedcb8ff940d159663baf203b367d50dc4adf2d3fb4ea2691f122d014860914b042603638cd5caaba56adff96718090f93a00f9f9423e4da66a9d1f12670def679e06fff10e4a5ea5ef064f9a4ef3c48b54504ff68516f28fae8eaf657a52db5f2b7789bcae3bba6eb21506add21f3cc83f29d84e6e31e9fd3fd07f9e367b958e8d3cbca9d37e8eea15f6086ac2d694dea8f72cc7e75ad15c7d8f599a85435e5c74c5932be6da572c15c34330ec0d1f5fcafef4c988a063b679108ad3117a5690290bdcb0c354c4b5e54eae6497c6b6889dfe6f3ecf6d8f319ac69b71187b37120ea404b349fbf61ede3f64580ee8f60a40e045e1a29fb1bdff436da28c29f22a5101363679f7b28963957a33e32cefd0224cc1d23bb8575765ac4e8fe66c501d8ebc8823eeb36143d59c8832d6349685270ee03161b5ab711497304e9f7e7c85e58cd63d695dd7bf10a8c9686b605d267302878757886023c9a42c15b2d862d4e409484589d4440c317c7771f4e18f6c522d046aa129a9b29472afcf09e26c72f3e4cc2b119d8496d115b311820e8cf5202cbe7ffb02169320c1b79a2bedb819ff012b81cfc72455ecd27a76b968a65ba86830f3438a7748961b9c8cdb68e319af6d1262b39572384d79118bb3095810abb43abbc2f13f9e095a9c97a9b4b8da1510c1d4b3c80171701b466071b5e74e6aabf7d5615352e5ccb625bfaa843fc3cb1f25e98d66a24dc0426d6ddf2354a0de5930daf758510d6af0d9dd596cc2a61a5e29cde73f559943c2ff01cc71f7ae49a7a6394e1e70b3fda6a3f5eef56b1fe9bff0c6193ada5135cbe8798273261dcf4c360df43ad1cbe472147c99295df4b65a3891de99c96197df913e9b6c431081fea7b1fc25d4469c4befd5906a3a7bb14535a2e29ced4c7f8a484259d5424c2acdd732009b91cac0e140c7ae080f4daa3986c2b67fbf45d526fcef7165796d9b36ca674f76625e500cae214705a1abfeb256658bb773776e943e1dde73e6aab50be3b6cf750427b09ce27c3061dc9d7c3113b329c90a60e1203d73831a830f808915415ad3d4cc8b729c7977d937116dc97067479f38e27b2517dee564307ee7713f596e1f095805ba4b19b3677b97486bdaaac83fcc40082d05f123ea8c9b87e0ba32f22421654f93710e35382419ea6d080b660222d16ec8b1ca81c375e43ebf4420f92e729446823dfd35408163b399244342386ba2e1af1da2c1e27b83d5e7f4d80c8d6ac79ab60d870b06347ffb00a59ff961f343afa9251cd9f3a42dd342334a2942f1eb1f9003bc59d5b010770631cfa50933f1164e75d5c07b47816df7a65303221960a22b182cb10eab69b5e7464054f4c51f07dbb829088c68c41a57226b6fef64deb29f97a9e41a80bb86d6d8292911ebc74f5244aded1939688f557b647566b8f667a30761e78d36f05df89f63ad8197b3e980a68a264a09e2a85263ec0ac9f7309412fa75155ad86680da6e0bf5bf1e8a069ac249dcad8f582a7acf50f0b8d11639c40e35adf09f6dbb880cdecafecc6ff173663e95cd6382e1b82d573ad0eef02b42f2a05df8a15f90f1b3dfb0e3d8886adba6377db11af3ca61c48130369a29f5057b9ab144f0145f15d714085eaa769221510e66bfdfe6df7d32688491751fc5020f4d05f12e49e59fb5a8679c79bddd9ff8bf4e38392e13f210272f661da20904b85fc0e3c22d950a0210f4a612bd86477e8dbd2eaa4b0398ab3fd6c095ddd94ea1325554d0c089f9e0e61e63c6fe080872f03914fe9d4f5c801f0a9ecc155c52b244ecee5c8f0cb66941f9ca18b053b85531ae882d4be998d9698e6261035ef8430f3621bed360dddff3f2001d797e59a260a0a527b64f61673783ec1b564f761e4101751062401f13648a2695fbc50511558f0b73e848d02ebaf0d90db28b84b4a7c8280946815a8f4ce4160ad4f2612af8cb9a7c168f754e551ca5a2cc93ae600e4aa3facac2ce39032f0e7da8f7691693a6742deb26dbb0653d7d80dc620de979e04a7717fbd0110ab5967a1819f19753e467820149cef06be9096acf2407d9add69481cb1a30d4c5d0303ab0609be4a3f3339da029c993a6948545cd0f111f93430cc261d60dca8424f242dfd9c5565eeaa2cd3e1c4f001ebccc694b294ac6189e9e7f487aa20ef2fed70c28d8200284ac5c11ae6f439cb2c8e50cac174bf78020b01aafa005979fac661c711bb3e7fb82ddc29da1d56f895201f581fe8712f91e8f01f52e1196947fcae943c4733963b510accb40775d0211ec64d9ab308e203142dafac54ae75014eb1a858ba876bb3982b3085ff1c696522f9a619f8de9ee7daac628305717c96d122dda9e759d5dbf642961e60d0c47caf9fc39581b213abdce39323efdb000c79ec2fadf56a7ad5721abf0365e6db043499a3b2e9fb6a4e36d527d55ed97500354b11cac5684e45acb2f9d3a73502c243950fe6b4f4357fef0a1a9877ab6ff658568a3a5e50b21a527541f412f4bc381abd90c37e44cdf5dc5b575b79b3206532db585925a44b8af5a4d974cda4a88fb0d45e1a2e0dfc249d09c1b7c35462c15c7b678594b156af9e019f9aab09053a22804da5d8d429497838289813e197d6d5440eee82d901f630cfa928a97a2ed477b88fd009cc242a0491697fbe109feb5cb8e81e03c714a385b6a34590233a6485b312a505506d88ba201ebd49700a527609006c647812e540553c148925d8ee6a89fa6378310feb3db7fed8ca7a6925da092b8f724c49d72f9396681ed88920b9a1bdddb208ee79d9f5715647f03a1f26e48564124167ce2f4ef60f0a78d2d887d1587d86b87456d69ef4b7e81f7374234b3cdcb1c3298e2d93c37a05a6f3eda9267666a5f20a89152300e8dfc22f406e6e35a3f04b8e53fa0476a4727ac31afb6539c1fe813c8d1ab5c458832edd35cc9cca2f80f22931edbc838f761a041cd97b3b83b06b493ca6d483157274e944af2a66fb4b468ec38a47d3e0b331ebfd173896db121d11313d9c49a2f6aaa28b15204edd2752bb35c9ef9eaa68948afdbd10264ba83052ac5adf34fac3d58943cabe8dc747f6fe90148e8d728d9e41b0c9ce6da0f9c6d7153e74ca627b3a0d7053c7495d1edc555870c20e5cd4a520829f4e78ac552915a07678871b330c52ebc8b846d006befc7aca9f03f20270a346833421d93df5409bd8f6a07152af455c35a3d95eb09353947dbecdc525c4aa4a24225eb09413e469ce2d8cfed37625fbb2f99277b2adae57949d9824e94c6847027843b38a4bcf829bc9c04af336439971b6d2c6f67d763759a16f61e1a58e478ddb5b376753d46da2e717bda24d421642ec1beb861c67fb1cea9e53024c3d2704168588205083bc13162a72c1596f5356835ac77b5a84b8a3b52dff9ec55374d72b158fe634a0a0ba0d064759c43b7347b38790267bf67192fb02f9d1eebb9bc3533a0a8d574f3596ceaee4027678a07c47be4532934965bfa3828bb3067a98ff05e9718b654e6b3163151c449b15620708ebb4c7dc6fa6fe44088a551f40cd3ecc3fc548d95116617c7974ecd2050700478591c823fdc07953a281e70aabbfed40b2a6fb26b66ba27bdddba72b720309d4cc6162007c051d5b9472876eba86ff2b3401b5fe24846a584c927710a725c8718c53ae2133ead4fede0d49ef9c4692aedaad6cdbab434c6c2e3b86764e0e68d9730d3d49b1fb956301b56699a2cedaff8f13f251f998a6b47fa02b462bbeaac16d3b4f159ec7515a75890ef4a082d64868a59ccdaa834db360819c835aa7bedf00d87f2af6cbce8f170dc46e9d03d28d00008299d3ff54be7eb0749467f44056441c46a2149c33a64fa22c671aa3a8b176d6639eab644f884864e901ffdfe2987fbe121e320b3824895dbcb0f510ab221930eff47210a21ca7530d8126eb86e28969e4e5cfd129c69c28d5400d6803251637cb2f4a9a35fc58807055c3e39c291a5fca68adb2887d569f7b7ac21114b3ecacc862dff3156ae8c975f08c9851cf882caad4f8f4395ac9edc78da7d2ecc1676de2428893e2630d47029af5bec86780fec0b59480461b952f7c9e24455af61ddfa42f6e58b0a27ffbb7e75f21e0baabe5cb78a48e182237c3d939cf7e2681c7f913d8b1c09a095ca804f9a63c5fde9fa42a373935d4011776b1effd053e743d1463216e126a50fbdb332e2db51e48e23b9b85ec0aefec248b555c01b208d39e4bf2691f673cadf75eee501939ae35afa4fcc23e8c7918faaa69b52f1b1b97b164bb626d1abde92ce55ecdc15b323027c14a2e0bfcb6b0214fa646bb7bc6faf7eaa0796982be288f45025a16f2deedceb88e487f33830d3d94013b1109375786a99a8e95809e484c83661e97b08ae2f1c762b3b4978798abd43130b751143e88ae35c98ace20f880947e2f5ca01fdf83e4475cfbc5aa42bd1b932be2e59b092f5848ad4771d8e41799f738ed6fcd50824011f18029d9772121db556ce61e531032661ab3b03183ef57efd9f74e083944320fd4753d727e0aeff15d01ef8ed95d45a82ffbf03e5f46e2cdbf54cb3cb3f4e217fac17012872df124bf0d81c9e99940feff7eeeb677aea6a90a699cafc92be186cc8634708f13914e07baa1f8c5bfa05fce8ae14d56bbc3d7ce170aa7339b568ad158cf52defc4f379cf47187e77b51756617df2c382096d41ed65205b70c6a0e43dadea317c92056f699af9d79f8d5fb3bc07802a9db02df086304f97a7552fcdc53d52f7ef4c5013d58058d65ffc61994facbf0d6894d4aa856c6595f503c6a9dcc2972316f20012962ca53bea672100d5332325a99e6e88afbe6778c228bb39077c43f692e8ebd2013447963f01cdc114fc4657a678421a77f558483df165bd087cfe10249bdfcf48134ab1686c4e934a070230c50024d55a34c657c67cd436c9ee9bc90dda0f7d71829f5fdfa5f8a80c346566ea748958c00297fa69d5faabcc0df8d76b18bab0fe04d63a43b36319dae700647391fda2615bd1c473231024c19a4f031b1e0b4edcfb3a694ac9a42cac2204086487db72467ab7364432ab89c5fb2ace069fcae4d2e3047085db6757333e2ffdc2409007926c0f6e1192eb77237f4587a728c36e2355bdab7b1e1752101558d85248e5a20c6f97380973a8d040640a1e43c4f48d6c75520872ef0bc158ff65f16c9924076344d05de2f0db943fc3340e338a4e12f6ecc296e14d6b2aa389387c0479a50b504d72394904fbffb608f3e6896a3e219a35c2548bb673661a7ae8ad3acb094910f678498edb1195b8b2ab4a3e6c4d03826f170cb0805384b443fb33498b3a367fcd203efb0404d8ee2cfbd182694ba1d39f5b0a22c39d596992c74e18a9a58db8840c0160a115e724c4373f7df0f826793d2662bf582937b6d4fbf01a96d4d57ca292773c9b3b5e968c2c01c33dad892d9e75c56852be846b387c6966c6e653cf131e2aba7d332378ce45aefaced81c8b6dfe6ca84e2d443cd09c2137e8ff58489ac1ffdf05154fe070e09d2e25d8aeae561389d1c9d04816df1d814a28009dbf9fb44a71f822a6e7f3e16d76591502380a60e06aebea87d44326c614e1c20170d109a6c3d2cc7581055c93f1516557dfe0c5eeabcdb410a5a90647b892f01dcf767388f34355f39924db4f1fe5aaddc7b10abb620c47a173da594708641bd540db3a58fe7a44dea56065c487d4a1232e96ed441b1131848d8a856e90c1495f6dd17d2f6a37b64f3bff3b63c28611a50ce715809d09f60f5f1238486ca7232f5b6290647f8ddd8fbb677797df4835eaf6c4bba76dd9431cdb5321fe9d1f16112d65d790dd00a77bc221354023436cd4219d3d3b57ece19efb8dc84333f76c8e6634042eccc07385a2e77d7028b5a66f0870ce1fcc6cd0ed6d432a27c8fc78caee7900c778d12f78edb73a6457a328a525788541255bafa4e5b5cfbbb0fe77ac619509307f0d851b8ce2be8d45b88770c42679d5a627914861478a56a5064ff967b6b51cf8b80bf81587eed18e98a018568b75f8166d41a696e539aedb54ea9cce2c9574096d98af486ea5faaf24a9e888921f47e6cf125c43731b38915e28d93afcc643a4332d9fe4c56e0822ce76928ad83f7fa9730d1459848a370efc44198b3b329f65b9565b885c62e28ee50e4bbdfe00140395a906cc91063c41"

	psn_privkey := "1f2317482fff7ec9b4137479289babe80dd0162a3fea3494cfe6951fcc6cfbd5"

	veri, err := v.ImportPrivKey(psn_privkey)
	if err != nil {
		fmt.Println("#1 err:", err)
	}

	ct, err := veri.Decrypt(r)
	if err != nil {
		fmt.Println("#2 err:", err)
	}
	fmt.Println("ct:", ct)
}

func createTradeData() {
	req := reqTradeCreate{
		To:        "1a9ec7b856e6a60603695a14eac6a96e9fae9501",
		Cid:       3,
		Value:     "100",
		Token:     "b526e1192ac5f8a5ceb3b2dd16ef888f2707d75d",
		Timestamp: time.Now().UTC().Unix(),
	}
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	sys_pubkey := "0464ce40c4e4340664ee346894f96a0cb3b5d84a6b18ddd37614ddb5e53645e1534d610fff186a1bb5852b5b0211773d40ccedce4a7468511566fe97066cc000b2"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}
func checkTradeData() {
	info := "d23feff33fd39df1a16deef1bea6457402ca0020378658aef39db91629854f0010022705127431c49a633f1d454efd7d76d67f9700207ea1d3e4f67465d462989c43aec98d3d216c1de49e68a6fcfdd71c359a6c44e9c28774d2a93ff50b02ddd6039d9db92f5b6088404fa1fa66a0ff295c7a3cd8199e94eb9eae9cb59b7d74f5a8b8d9ed25d52acf1363e339e107371cb89b90579b36c44ff656c2bba29fdc4a12a001b7f5"
	req := reqTradeCheck{
		To:        "1a9ec7b856e6a60603695a14eac6a96e9fae9501",
		Info:      info,
		Token:     "c3322ab5d9bae0daaec69ada70fa2547a6ff587e",
		Timestamp: time.Now().UTC().Unix(),
	}
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	sys_pubkey := "04a017ac0e9665ed7dd02df71ace98fde3b2a6a109f4181c0e397392a5e353d36d5b6d7301fe24516058ee5bca4172dfa0bdfa78b98fd24eb4a58da0441f665731"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}
func confirmTradeData() {
	req := reqTradeConfirm{
		Tid:       "5bca3643433a50e13d6c9e1826ae253444f1f5e7",
		Token:     "c3322ab5d9bae0daaec69ada70fa2547a6ff587e",
		Timestamp: time.Now().UTC().Unix(),
	}

	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	sys_pubkey := "04a017ac0e9665ed7dd02df71ace98fde3b2a6a109f4181c0e397392a5e353d36d5b6d7301fe24516058ee5bca4172dfa0bdfa78b98fd24eb4a58da0441f665731"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}

func confirmTransferData() {
	req := reqTransferConfirm{
		To:        "1a9ec7b856e6a60603695a14eac6a96e9fae9501",
		Cid:       1,
		Value:     "0.9999",
		Token:     "c3322ab5d9bae0daaec69ada70fa2547a6ff587e",
		Timestamp: time.Now().UTC().Unix(),
	}
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	sys_pubkey := "04a017ac0e9665ed7dd02df71ace98fde3b2a6a109f4181c0e397392a5e353d36d5b6d7301fe24516058ee5bca4172dfa0bdfa78b98fd24eb4a58da0441f665731"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}

func createOfflineInfo() (string, error) {
	req := reqOfflineCreate{
		To:        "e7a5f1cd76e372a6dfdbfd404a9e573fe6f48017",
		Cid:       1,
		Value:     "10",
		Token:     "c3322ab5d9bae0daaec69ada70fa2547a6ff587e",
		Timestamp: time.Now().UTC().Unix(),
	}
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	sys_pubkey := "04a017ac0e9665ed7dd02df71ace98fde3b2a6a109f4181c0e397392a5e353d36d5b6d7301fe24516058ee5bca4172dfa0bdfa78b98fd24eb4a58da0441f665731"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}

	return pt, nil
}
func checkOfflineData() {
	info, err := createOfflineInfo()
	if err != nil {
		fmt.Println("err:", err)
	}

	req := reqOfflineCheck{
		From:      "e7a5f1cd76e372a6dfdbfd404a9e573fe6f48017",
		Info:      info,
		Token:     "b526e1192ac5f8a5ceb3b2dd16ef888f2707d75d",
		Timestamp: time.Now().UTC().Unix(),
	}

	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	sys_pubkey := "0464ce40c4e4340664ee346894f96a0cb3b5d84a6b18ddd37614ddb5e53645e1534d610fff186a1bb5852b5b0211773d40ccedce4a7468511566fe97066cc000b2"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}
func confirmOfflineData() {
	req := reqTradeConfirm{
		Tid:       "faaa0e39fc69eebdbd56b32cbd1f1c8c346e8968",
		Token:     "b526e1192ac5f8a5ceb3b2dd16ef888f2707d75d",
		Timestamp: time.Now().UTC().Unix(),
	}

	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	sys_pubkey := "0464ce40c4e4340664ee346894f96a0cb3b5d84a6b18ddd37614ddb5e53645e1534d610fff186a1bb5852b5b0211773d40ccedce4a7468511566fe97066cc000b2"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}

func createNolimitData() {
	req := reqNolimitCreate{
		Cid:       1,
		Value:     "0.1",
		Token:     "c3322ab5d9bae0daaec69ada70fa2547a6ff587e",
		Timestamp: time.Now().UTC().Unix(),
	}
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	sys_pubkey := "04a017ac0e9665ed7dd02df71ace98fde3b2a6a109f4181c0e397392a5e353d36d5b6d7301fe24516058ee5bca4172dfa0bdfa78b98fd24eb4a58da0441f665731"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}
func confirmNolimitData() {
	req := reqNolimitConfirm{
		From:      "e7a5f1cd76e372a6dfdbfd404a9e573fe6f48017",
		Info:      "174828dad146f746f85b10df373a147f02ca0020fc749098266558c5724c9a9fd82525598cd0741690ef188ccfeeaf17f91ee66e00209d6aa24003cec896f21bcab9975764be7787dae74c850c3bdd36c7386eefcbb8f3f34a7f0927791f2d0b3d46ff4e45ca5317af89089f10443b42bb73acdac81b7d1c18f61b10a11cb8eba8cc83f34be436c0364dbff05a41575457f63fe1f7e370a36669e8b24cb36bead33f5f0bb017",
		Token:     "b526e1192ac5f8a5ceb3b2dd16ef888f2707d75d",
		Timestamp: time.Now().UTC().Unix(),
	}

	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	sys_pubkey := "0464ce40c4e4340664ee346894f96a0cb3b5d84a6b18ddd37614ddb5e53645e1534d610fff186a1bb5852b5b0211773d40ccedce4a7468511566fe97066cc000b2"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)
}

func base64encode() {
	// data := "5b21c025daaef7e877c5440695dd308c02ca00208db3e52513333d138fc40245c844b5f02ec80167ae0386b365a63ec214e68d060020fdd5f1918d01e31fdd1175992c95ef6804b78433a363ddfc286695e9256fc807fa729c3259a5adfec8180c1a81f6f884b63a144555b8fe8de1d47bc473150c87c7bf9fced68e85ee0e84bf529367c5c93d8e90d8a0215bcdfb2827a14e37bbc5054c0aa9156d3f2c1d531f1600f339a5ebbcffaf11b852190b789f2f5f556afe"
	// sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	// fmt.Println("base64:", sEnc)

	req := reqTradeCreate{
		To:        "1a9ec7b856e6a60603695a14eac6a96e9fae9501",
		Cid:       1,
		Value:     "0.0000000001",
		Token:     "b526e1192ac5f8a5ceb3b2dd16ef888f2707d75d",
		Timestamp: time.Now().UTC().Unix(),
	}
	b, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err:", err)
	}

	pubkey := "0415216f9a6111995d9bd7db0d9ffca5fdad734595d74f0073392570018be21d32aa1de68337f44c62772ff7ecfc713ccba443314b23dc44206fe758404e97d107"

	veri, err := v.ImportPubKey(pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.EncryptBase64(string(b))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("encode:", pt)

}

func toLogin() {
	account := "0xa83c7f79f5BF4C0Cc883E44ca9c0987BdD94B7a1"
	sys_pubkey := "04e213164933984dbec638099abb368ca1a774ac19b89fa96677c5c2f196eddb0851643c8a0b8c21f5be0289d352d1e2abf2b591749ce8d8a1a74d16cbf1ec10e0"

	veri, err := v.ImportPubKey(sys_pubkey)
	if err != nil {
		fmt.Println("err:", err)
	}
	pt, err := veri.Encrypt(account)
	if err != nil {
		fmt.Println("err:", err)
	}

	fmt.Println("encode:", pt)
}
func backLogin() {
	r := "4f9ae69f07fe72046777e336f1a1f3be02ca0020fb83f9d64e233ac068ad0dfde422f37de21e70ffeee3c91f6a127c887bf47d3a00208e81af5c9a575164c0862e1ae59c77cf065807310ecf393a141f8505ee6b78e83e3695a7849348727704790f32cd4c414f13ba96165f17ccb68c852ca36388874ef580abad05986a040b9f72d0dd5a4c89f50bef0d44a818a31a203bf02817d7dc63b0035ad388abed9873342af9637b"
	privkey := "1f2317482fff7ec9b4137479289babe80dd0162a3fea3494cfe6951fcc6cfbd5"

	veri, err := v.ImportPrivKey(privkey)
	if err != nil {
		fmt.Println("err:", err)
	}

	ct, err := veri.Decrypt(r)
	if err != nil {
		fmt.Println("err:", err)
	}

	fmt.Println("ct:", ct)
}

func calcWid() {
	timestamp := util.Timestamp()
	account := "0x1F6998A6153378aeC050Ca02ec37B2cFC6ddDbF7"
	uid := "9f5da975cc1da468f73f01624acd043768819648"

	s := util.CreateSHA1Hash(account + "&" + uid)
	wid := util.CreateSHA1Hash(s + "&" + strconv.FormatInt(timestamp, 10))

	fmt.Println("wid:", wid)
	fmt.Println("timestamp:", timestamp)
}
