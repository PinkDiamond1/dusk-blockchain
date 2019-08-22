package config

import (
	"bytes"
	"encoding/hex"
	"time"

	"github.com/dusk-network/dusk-blockchain/pkg/core/block"
)

// A signle point of constants definition
const (
	// GeneratorReward is the amount of Block generator default reward
	// TODO: TBD
	GeneratorReward = 50 * DUSK

	// ConsensusTimeOut is the time out for consensus step timers.
	ConsensusTimeOut = 5 * time.Second

	// DUSK is one whole unit of DUSK.
	DUSK = uint64(100000000)

	MinFee = int64(100)

	// GenesisBlockBlob represents the genesis block bytes in hexadecimal format
	// It's recommended to be regenerated with generation.GenerateGensisBlock() API
	TestNetGenesisBlob = "000000000000000000e2d92d5d000000000000000000000000000000000000000000000000000000000000000000000000a67cf863083e3e4512ac3f697ab16754c9fb0e9a21515c7982851177a79d9e73c805fd051bdc80a17a00f2a8604089bb8c828a6c5ede4fe32498e7f6802d5af4a60000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000138780b7acf5c25eab8cb3fd5f095c3acfc75bab397869b787cc4c5a8adcf6fb3400e0bc22d7accc8a1de91aa162f71bb225f41db495db444a0d8f844192b604605f87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0a0000000000000000000000000000000000000000000000000000000000000000bafb394b549d97e178b901ef210e15e1f47eb573897edad59ba999d937707e7c200040f09bbce1080000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000c6b78d81fde86579aea5ced173e9467a47d0e3aeef7ef6428667a6531c9bdf72200040f09bbce1080000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000c6a76486e45b9a6493376952ab66aec97883c1dd9de5e123c63c47e99b06117b200040f09bbce108000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000042e6a41d518e15ee69fc2d78dc3b17daf2fbd356d9092ac18f945800cacbd911200040f09bbce108000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000016ea9bcbb1e1197c712901beb52304c0bb7c3c811b98fe5c189eaf67fd57c249200040f09bbce10800000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000002219e47caa2eb761babbf774a07116ac5c6f1b6523e5539436dcd0bde4ddf900200040f09bbce1080000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000f084fe363e335496476628912abc04d45d3ec310e4a5c37b2593fc035dd9d043200040f09bbce1080000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000beda995e37346db55af70298fb88dd7804c5004f6a2523710bc4dee39049d20b200040f09bbce108000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000020faffa5818c8e93bee71a48985fcfd5a79d84b65e25208714b65779d6401259200040f09bbce10800000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000005a0a891f6394559434eb56d427b146fc8db4762a20b31f1a194861c14e7bc103200040f09bbce10800000000000000000000000000000000000000000000000000010000d42e5adff8f41660cceb3c9b79e76d5fc8e4c75b11cb398ec6c7f3e80eb68f6987d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000001ea262f17332ef2d41c5872dffcba0e5a09757d6107e52aa6113166251e30e4c200040f09bbce10800000000000000000000000000000000000000000000000000010000186f6ecf955bbc4ba7992967153601773c06d4bac287160e80c2333345189c7c87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000007035beee16ebf14975c387275915d756bee008392f4dd6accb7115660f572b3f200040f09bbce108000000000000000000000000000000000000000000000000000100007ad0f993c5b5f29ade497b7c8d626d52952f419c8768baca159f33cd892e8e3987d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000004c4d99658dfef00bab5cc0e5a8f4408a7098851f601335950b5380910db3304b200040f09bbce10800000000000000000000000000000000000000000000000000010000406717a9b1889d33f2786e69a7e3252d1735629d12187f64293daa67f542f56e87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000629110dac99c896225ed093dc25cfc89aee4e510f57872182a0ababeb6b2b821200040f09bbce10800000000000000000000000000000000000000000000000000010000406a61624ee4299f7b99c7af9caeb39115328ed65870bb4f08d8cf34be83d81587d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000a020462ce90d8ffabe91e3b2ec04cbc017dee9261aa996ac6a00ca958e710c55200040f09bbce1080000000000000000000000000000000000000000000000000001000092a90aa8aabd647beecaddf690130fa6ec61ae12f15a37c4dd660cc83ff4f57a87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000004c82b4a97b2e52014ea442c47aa381a668905ac2bdcfda722011d1b8b6d46068200040f09bbce10800000000000000000000000000000000000000000000000000010000720d33c6803cd5e90f418d63ffc16f3283d9feb3fa77366fa0738844fec3f50e87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000028f6e268cd96ae3852dd69bee881ee14269890ca121c3cfb49984e8dc5f8356200040f09bbce1080000000000000000000000000000000000000000000000000001000004bdba6039d3f25704c8cb7c0da6c052e6f519c7c792cd3e9256dc69a6e5ab1587d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000bc6036c5c74b1b948b116b7bed612009dd1bed6a7e969504a2c9e8a7d0e31913200040f09bbce10800000000000000000000000000000000000000000000000000010000a8c607cfd732fc472332de95efa911666600c53655585a38fc88f14c1ff0192587d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000fa46da03bfc5bafe59a9c8ba8803ae2682e2682c0f4e2038e1813afe5d911c56200040f09bbce1080000000000000000000000000000000000000000000000000001000074157dd6bcc5e51d2ae237c40ca12c426cbc9d5925e0f5401fb041fed391974987d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000002098d0b51b30405e918a0215f58ce3ebc107cb9f3724d3be01b76db2a17a7b67200040f09bbce10800000000000000000000000000000000000000000000000000010000c2e7a871b0a9905d95c7336125cf0412edc8557eda4ac53feef631c8d553224987d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d01000000000000000000000000000000000000000000000000000000000000000000793563c0715a35ecf732326eb178f7a551b8fcfe71449d9885522a161f995f200040f09bbce108000000000000000000000000000000000000000000000000000100007a6d11eee0232f698f26181c17838e66d2b58a19badb3bd4d423bd3f32dba40e87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d01000000000000000000000000000000000000000000000000000000000000000048b6983fcaa8813fa0ea7fc7dd032ee8ed67e674438422b1b1a3a89225717e18200040f09bbce1080000000000000000000000000000000000000000000000000001000064d11620914a420eb78f3306fbf4b2193321324e89c95ee37098e6c11a21e31687d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000001e7b26c1993c04a8b2e918ce4dd673d233cb67236bd8318775cffb70eb07a237200040f09bbce1080000000000000000000000000000000000000000000000000001000078fb95f6a538d7a183e3c9b977ac5986986d33585e319ed8923e43ef4fca137087d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000985eefa7ff432ea7df1ab27293c2c9dd7599f997bb929fb0e8051990b9678b00200040f09bbce10800000000000000000000000000000000000000000000000000010000a2eddf0a06e61758660ec973e4f8b33e9f6b73c95d9165d9c4bef6ab22fe7e2387d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000007cb03a689d5e22642b781830f48648c82bcba661dd4a24c1e60c1c921e320849200040f09bbce10800000000000000000000000000000000000000000000000000010000ca624d2d17a12031dd1f00b62066de9970631464b501425389450cab11317d5e87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000b00ef40aadbd88c489e11688d567edcddf0597475e8d99be757437261d615604200040f09bbce10800000000000000000000000000000000000000000000000000010000fce4839328cce60b6a86fbb1023b162be724384c23434bf6eefaa0a83d66ed1787d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000001c395cff6c0e4bc2d542bf5ea13e060c9d05a16f3558b9b68ddec6a521bb5276200040f09bbce108000000000000000000000000000000000000000000000000000100009ce083377fd2704f4357501991436916eab1e6048a20a2062806df2abedab76e87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000382fa3c8760a3dd60d8532e5cdefa26c65237391f692b2223282b1eb8d415c45200040f09bbce10800000000000000000000000000000000000000000000000000010000e22fe5977a788fa1edefa8bb7628505e9adb2ec306439749aed71a3b471b2f4987d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000065ae192d5d3656134269ab7eaf71486a27e36b77e93f78e50732f3353199a36200040f09bbce108000000000000000000000000000000000000000000000000000100006ee185c731f0e44575746f6a2bfee2c2f22d46245322a1d18c68b3b2853fa40f87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000008c9abbbf0858027bf73c954a629e1aadf47ee09defe5d26d2412fc8d49be5f11200040f09bbce108000000000000000000000000000000000000000000000000000100004c20d8fb2cb109ad59f2216cd2c98916abb11a074c70ed7b8ce87da29dee393587d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000882172ce382626977c0d620da216d2ca7897bec1895522e9ef0e37971b8e122a200040f09bbce10800000000000000000000000000000000000000000000000000010000b6a7f8c8fb78e722c2b18b7520397ff2b834c9241fe99e6b5009a864e763052d87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000f0422c80cbe1afbbdc452bf0abae27522c94a9383c8e4105d4d801947e814344200040f09bbce10800000000000000000000000000000000000000000000000000010000bc1282f539b3e714de451485e8375c0d9f3310226dc42f72382cd71dcc61b22187d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d01000000000000000000000000000000000000000000000000000000000000000002c0c6256e56ebd2e891ba7c5ed9d81743d3c5b597ba3577c2ab92b604385361200040f09bbce108000000000000000000000000000000000000000000000000000100005af0a092f3c55cfb0f42af2e27c39974e0440f8f772cc339e955582b63a1397b87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000005aeef815d66cfe98c38a35078a4b03f96fb22212f13b06482256fb14e15f6606200040f09bbce108000000000000000000000000000000000000000000000000000100004a67cac9c43b2998756e1c66b796da336eff8ec8f7e62092f56fbbbea736d86087d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000d2ed827884f10b775aff1dfd02b5058bbd9a0c823ca400ef43e066a048d8371e200040f09bbce1080000000000000000000000000000000000000000000000000001000020be1361906ea95836b1b8f7df2d3c2d165dbb39fd78e67ceb97d92299313d1f87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000004856f73898bc15f82b546e0e54329a5302d8fb450e324b5c412e9b1f75fb4e13200040f09bbce108000000000000000000000000000000000000000000000000000100001cdd0f606b863a59d4e934dcc6f6a7f2bc75d21fa81aa1ee1da0feb1a4a2065987d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d01000000000000000000000000000000000000000000000000000000000000000092a36867ab1ec455ad0e0c8a9fd4f4eb8dc0fa003249c08bf5f2f00f9d558503200040f09bbce108000000000000000000000000000000000000000000000000000100003401ddd4e70e0e5b48bc9a7cbc285e6b9d889452a287cee3037dfcfa34a3a42f87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000f2678aba3d340027ecb38834535bf686c0f91c0e328c64c4368164d1e097471e200040f09bbce10800000000000000000000000000000000000000000000000000010000a43ae9d46447d69748377efba061a0e823a7bb168ebf350dfe2c771a28a4211c87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d01000000000000000000000000000000000000000000000000000000000000000064051de2be0fd726f5873fbfe71b6ca545a54cb111959d7eb0346102cf5d077c200040f09bbce108000000000000000000000000000000000000000000000000000100000a240d5847ea6de56082b81554aba119f3ccafb550e8926487c0a049e1a28a0587d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000003ad901f5cc28cef58832869517d300c39336bf0000b5ef39bd7d88669fdbda6b200040f09bbce108000000000000000000000000000000000000000000000000000100002a53c1b0bad671c92b714b39279535e62c2acf9ef605da6c904cc22c12f9f56d87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000de407f197e218ebeef808f53e7009a66aa80c0f8ce7680eb99d9fd209150ba05200040f09bbce10800000000000000000000000000000000000000000000000000010000b0cba1c0ee36330742989df918bb3e1932db0f61ba2f8b46ca63ae4579c0b97d87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000008c53cf0dc74db1d2fa6d28ad9ca133ab72473f7241f33ace88f64e4781a96522200040f09bbce10800000000000000000000000000000000000000000000000000010000fa8af020924452e2df84e2e95183648d8cc1a41bf28d00f903d1d52d1ec8407387d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000eae761ddb12d9e4ac4df295184d74d83f03b8f0725f6cd9e56acb9f4cc0e632a200040f09bbce10800000000000000000000000000000000000000000000000000010000462d058c97e490097dcb74fd8ff68d7fb1cccbd07a0e1795b19e2c1ea4ff622587d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000004864805ab00791ca27173f9b8567636793d07d510a23d86eccf186316446e933200040f09bbce1080000000000000000000000000000000000000000000000000001000062d8bd208e9dbe88df00f6eeb2760ace93f2d69393d490d540687e62997be61887d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000a685d212020705fdd928f8d161a781ac834dabb76d202094f8eecf859baaa445200040f09bbce108000000000000000000000000000000000000000000000000000100007e76828aaad5eb1b899b87f535f59f40fe2ae2e609b4c229cf7d197096acd81487d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000ba704c232ee13ab454eb45f8a659429bdc1fbe6a2d5b9c2a6df918827e0cad26200040f09bbce108000000000000000000000000000000000000000000000000000100004c649931d14fada8f14c763c219963b9902742ef25bdb06e26e3a6ea946c171187d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000664ad20f32cde5d0ebb52d1035a8825f5c611babc30ffd3142eaaa44aed69554200040f09bbce10800000000000000000000000000000000000000000000000000010000a2330e37f45f1e529287dacc135a5806bf34c9c1e2f4c94406d1f39ea95eb14c87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000009caffe4da24bed40710b209c36b4c6197785f4dc7ec82e154ad8525025e8e004200040f09bbce10800000000000000000000000000000000000000000000000000010000ee56fddc96f879accd3a7116479e40a0c7c37f072885d2e7c1c8b6063554336087d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000804d86170b8999e9007154d3d0a72482a8c9df86c5e735dd8e4cbbccfe494823200040f09bbce10800000000000000000000000000000000000000000000000000010000f82ec9e8595435bc96d7ceb26e63269099a26178ca4a4ac3ace260849d6a001d87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000005e15df43d66e7f21c354ed9aa2010429dcb2cf2d9f01c2323be69fff87431173200040f09bbce108000000000000000000000000000000000000000000000000000100006890b0f055a757ddc2bb5af3b5b3f28bac0761352395d9b0cd014f1287e6907b87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000def676d56233f3b14b8d978cd06d28e8223a35abcf465e07ae45705320171e33200040f09bbce108000000000000000000000000000000000000000000000000000100001227950ab9fb2f4969486e304906f0c9d72079875fe6385c6af381bafe519b1e87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000e4350c22bfb82dce3f5a5591c507c806a42479bafa948dcc55302b067adda62b200040f09bbce108000000000000000000000000000000000000000000000000000100006a7e7769420585aa3deaad6408c7be3d4785fc2d2eab4f0c709ed0303e93cc4387d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000cab6692d21170d3e64f556505cd883b8a5bb9a65dae400af1607952f9edba172200040f09bbce10800000000000000000000000000000000000000000000000000010000f609fedbf7ade155b5282b848e67bef402ce8d564ea18210255dd605d5b5335d87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000768dd42b60eb19b0854cc212415b7660f9b9a10321e718a463783800207e305d200040f09bbce1080000000000000000000000000000000000000000000000000001000014dfefee852023162ce1bc960904c84a6d07af8154b43a358e2fe3e6d278e65c87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000007239792b497f3ba2c65f87fe0eb2e2d3dfe2a686d2d826a72ff7c09a88eee403200040f09bbce1080000000000000000000000000000000000000000000000000001000098030c0bf2e04db60c56eff856601390d71da9a6a132af0b0c781cc0ef305f2787d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000007466fad30adf20ae1ca89f1d5e2136d34d0061b39dce35f1c489e10578d0b44b200040f09bbce10800000000000000000000000000000000000000000000000000010000bc85d8b087ad03ee6c18646eb36f62b5d0c07440acafa9333980f1800b9f0e0287d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d010000000000000000000000000000000000000000000000000000000000000000bcd32d4effe0c241e918ee939a372bc7b5548fa6f401d2f4cd818bad799b242a200040f09bbce10800000000000000000000000000000000000000000000000000010000ec52f2f47e69758c744e3d4625e8d4a61578bdaeca04d6cabb7c5c4c26e6b41b87d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d01000000000000000000000000000000000000000000000000000000000000000008c6364849a541d5f7a76b9e25db7c84034d90f04c5c427aaecfc3d23185e35a200040f09bbce108000000000000000000000000000000000000000000000000000100003489047bc35783ba238edb8c1e80dd0c9683edd3e94902b3d0ed694624daa57487d1929f4c070e86d4d569f71fdbdb6a526ef47fb42d8ac7cffef4bfad7f77eb202ccb3ae06f0c5e66b31a095efb8262037b00cb547216bcd633d2833042ac214d0100000000000000000000000000000000000000000000000000000000000000007e6d909beb114ba8bbfc01d3ec816e744b8fe3bd83b4747ce4399fd76c506c10200040f09bbce10800000000000000000000000000000000000000000000000000010002deb57abc541c7f8b6ab38bce6082d1b1ec1105bf5ddd1e45cd4a17ac2d64070b000176e4d335b0f03db4336eb53bbe795cff1ec9c5baf10c7e4a1dc02548ff1d556116ea9bcbb1e1197c712901beb52304c0bb7c3c811b98fe5c189eaf67fd57c2492e642d6ecbec6dfcce363d0b8e45929bf0c8fe569833a3d8f25c719d5e70251ffd280406acbe55234f06b835a00405b79aa7be95d4a585aa359988021aa37e7d47c20c0000000800000002bae946ca42f50fd8d4970ad588d859cb89672af5c33dc3b6d74a2a83eea55405d41f9bf42bad2ece955e45e092562691b191dafb48720d79ecc11b0ee6e0e30b35458557e484083daef4b7c56585e4dd29e63de7c6d203ebed63f4a9f87dc10336818bc8d495cbcb413086f967050b79e196b80b4010b11e11e0d6ab6432b90c204e38a33fc22083a8611cedf8f763a0510c4a30b146394f2e49806693084b046a57eb8e57105d17bd55bb0dddaaabd8bfe2051aeef05911bdec0be0aaed4100d9feb06319b09a55b5f7da3bbb8364b8225daa17d1638514181631c147606c0668e035028f42bb07bdd6710f5b6f908819bb3073d651cf3659b21b8cc08c0e0c42e0797431526c899eeccdfe15bc50dc93fd2641149d76693a4c196fd2f91e04a8d9cf4a34a1dffd5ae22d017980a7102f845e718be4d4d6493cf8ecbde8e50ebb27abe9e9c71eef4074cde5b0759a76e61700a72db14bf1fd43817301adc30911abd7fd8e685610bd0eb91fad22732ffc7ea0f3fb68c7f48f74a08b5ee72f0f370151885a9727f80cb3ea0a089515f85757d3f5fd74c342b6c1b8f5417cd003bb1e68f2094b852971c06d4a856f698002b229613a46e1a1847550d1ce1ccb0a9d3ea75e790f448c258b3ddf059ef3211535b52b09421a12c58d37886d140f0cd88a0487c12782f0dcae512bac439722ad2c09d9c2d10ede2446dfb294166d0708c6364849a541d5f7a76b9e25db7c84034d90f04c5c427aaecfc3d23185e35a9417f0000c65c64416fdc796f5d571479fe41cea3422fc0234cf3d72d85fd34a00793563c0715a35ecf732326eb178f7a551b8fcfe71449d9885522a161f995f661b09ca02cbab439a764390cb5aa69de97cfc413857c02a5bc8edef02ed77431c395cff6c0e4bc2d542bf5ea13e060c9d05a16f3558b9b68ddec6a521bb5276409ebae50dd0039aea44b6fc5a8870ab1fb4ef67e8e1bb76c99d30293bf0832016ea9bcbb1e1197c712901beb52304c0bb7c3c811b98fe5c189eaf67fd57c249ecffdb3940e82d77a06936678fce605c2766c16d2b557abf9f6bdf866b95d23016ea9bcbb1e1197c712901beb52304c0bb7c3c811b98fe5c189eaf67fd57c24966d8809357c85f6e2e0f8381e9c89b5aeb5f1206c7c30c36336cac199ab39c3602c0c6256e56ebd2e891ba7c5ed9d81743d3c5b597ba3577c2ab92b604385361aa0527c39c95c9cc3565000ae7b12204f2035952b8cd7d7135ec156b0b5fda16028f6e268cd96ae3852dd69bee881ee14269890ca121c3cfb49984e8dc5f835644e64afb1766a37223febc1fad507db47502cfef4da4ca169d31e7a82ce07465065ae192d5d3656134269ab7eaf71486a27e36b77e93f78e50732f3353199a361c9854af1298016abbfc146d71f8a017f3a4efedfc62fa934e747568615a275002d67a9e0aaf48410cbe8c69a86a54808167a7dcb6c471ed7bb8033b37a7072a10f0b9e3a4dba28f12efec34fe805141dfabcf7e3abd20426d63264e0983c9384b2000e40b54020000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000002e89dd5ec9333f1e3df929018099a61605cc0a49eeb616e11f50ef7127595e057e9805b6c93ba90282ad221e595e6af1334f9db141233697600fe0338d041f4b209c5be447bae108000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000006400000000000000fd240300000002d67a9e0aaf48410cbe8c69a86a54808167a7dcb6c471ed7bb8033b37a7072a102e89dd5ec9333f1e3df929018099a61605cc0a49eeb616e11f50ef7127595e054c92d9cb759402153e161fd50d88a858927e8126793a0cde498419c76c83e17866868b9d598fd0f7296e00edeb2b8a46d3d761d79a6e33e3bba41281f758774d2c58c15dcacc5580d95b93dd229e859c30946e935c1fd5d2f2d2db6fb315017f3029f5268bba9105f1aed730900959b61c13a953131a7bc8f473a88e5fc908352b998e5be1a402ca3433211c4194e7eb51af99f7ceb44f3b0deafa969022cd0252f216ee189c35c0e59c874fb7f8b12cd74f8219256f9df3285e9edbb0e89a01e217c4c1bb65e14347aa9f3311acbb2471eb96626bfcafc97e3bdcaaef63fa0e884d5b1bf25578d89db7f96aad8f386a864910cec7a85c0b0a22d06e10c3030f587ef1bfac50866e013c66002c5f7d3ea4f500a207ac6ae423b06f051be3d900ec86b08825a5e5761149c84f4bf8f92748400b7303a57efac62263e526c4804690388b83e96cca574d36a96a50d43f65710b47518fa7a3b953b2351f537a0d4d78d9ffdfc22d28beadb0cea34641e817ad6a259131553c163d76a81f2f7cd17f1680a1fad6a548c5723e493247a2a8784a6df3e85f4b85b0f9c4dcec7f864e21ce04e99f9fe3bdd75bbb0050867a7cf74569e579693299069af5e51eb76023554455dfd9aa14779574c3a030ddab583c32c70356f2feb0222b4f3d9ee3ab4e6cf8c1dcb4931ddb78a333a6de5f134ff3d7b55f4bc0c73bcc1fb8dc691a0327761ac5212c31fc6f8318806b796f091e43d2cea5bfcdfe29faef5954c6d5894f33e447d7cae1809613947ec3ecacd1df1da92fcdb9d7059210e2dcd69ace324f6e8a870ff424cb2d48158d19433ef0b5a82f807f3abd2ab7042731ffd266f3a100e091d92b71dbc22259c34ca42acffdc41ee15541e528d103bc8d66855997823b1e306bcc45515e2c01d21d4670670fd7b043665659a4d005e439c814413235371aa41bb4e034423548a4a553009a71345cbf9fb13dfbd4b4f9b53ac9ecec8825529db14e82b0e302667cce23cfe726b187fd2ed482a72fb58571db345d736e6f90d0030000000000e9fd3b62d35c7e451cd2329384a38827752a0e7498f33f6b9560dfd80c7b7bba8101660a6b1748fe773681e6cf7cab008a169b901cb3350cad041a509f8ec1b5d0fa005749878cd2778fc4d9cd3c6e1df252900a6fc484a29e33670945d596876e058e069b7f4e751fee1f215613007339cb09e33e54538f8adabe8c466317531d508775cecfea6c32411b8a201645286a605e5e48e233cf543103ff48362d5f200b01aa5a89312b0bb06efa549f4aa921d6a60e09dd559e9b1d584f3b9d63fbc6f24e000176e4d335b0f03db4336eb53bbe795cff1ec9c5baf10c7e4a1dc02548ff1d556116ea9bcbb1e1197c712901beb52304c0bb7c3c811b98fe5c189eaf67fd57c249b4d025396b5d37a37c05463fe8924dd6b7192bb01d2eae3d41853d958c71bb12fd28044c7b1a661369d6d15256093e18f421ba5999460855394a6c78f4b719868b2d0d000000080000000203cbc729ac80be4e85616fb32ac1256e32492ed2b586df7572283ab480144a0d4807447a16d39a0fc0afd94a5eaa447f7f751dbc12035989188235895497c50907376f4d085ffc972d61899c2796d31a57272bbe282dd681afb600aa8afa61070a0e63b7ab355bd3898b6e58a7e4411e3d4c06b817cdc2954adcb1a70b6ce20048e82eada278cc53f3cd9888cb7c48e9076a01ce4dbad85c1269c5f97d5dcf0f7b50298800ab343e6c67982ac4d53a8341bc85cd593095a98405ecee3b86eb06bdb2dcd37fc7e72d8ae8f9e264bf98da14345e7b819c9c983ed298c3e29b8402e2d8911b631019156d5a421dddd462730caa6946efffb481487fbd01dffc460c0492147292efe90e3d1fb2e1ac9fc306649e378e42460380bdade7a19185010e656a31ec2ef9df46101c25163ea9b1a30057a4f73e4badd221989a69068d880cec7976a288e8b8e2e9715e06dc6a55f20cd22a46facb8245074f2e95af5b8704eadeac7c88cf97b788326df307cb9b7a234f1f5ae29094a7421d3b704a07ec04a4ca4d034e61dc114e9d935e963d34458301ae841d9ffae37bb39ab593091b0231a3e2b7c8b0baa44cfdc3d5e2e80025147a9df40f890d34ee84a696a5a5070b5041da61083d8eb0b1b45446d3f21e7a3b4a5b9515906c1168e10bf2acd11f01f749baa2f25eb02c24269ba39b71a66bdefab2ecf1b9849f10fabc6e14a7030b08c6364849a541d5f7a76b9e25db7c84034d90f04c5c427aaecfc3d23185e35a601d0c67ad793dd196716e6e5efc375cac6e93dc425b0b753182ac96c7a80a0a00793563c0715a35ecf732326eb178f7a551b8fcfe71449d9885522a161f995f0635451d683323cda41c49d47516617b80201777d057bde0769d2f502d40c8621c395cff6c0e4bc2d542bf5ea13e060c9d05a16f3558b9b68ddec6a521bb5276ac5c60a8d0377b3df567fc8985d0b8442786f2f202969dfd8c3e41a3b1dfc73516ea9bcbb1e1197c712901beb52304c0bb7c3c811b98fe5c189eaf67fd57c24952c0589cccacb16fe739ec1f11a99bf6cee5d591e8a2f0003961856ccb9a717c16ea9bcbb1e1197c712901beb52304c0bb7c3c811b98fe5c189eaf67fd57c249428963a449194020de8df5df05e88fa3b6d4cfa77fd547fd6f98d4a59652252a02c0c6256e56ebd2e891ba7c5ed9d81743d3c5b597ba3577c2ab92b6043853610ced7b8dcd261481cd3ff6bafc8e3c6dc3873a8c73a63f416f1dacb4f99ed43b028f6e268cd96ae3852dd69bee881ee14269890ca121c3cfb49984e8dc5f83568a516aa2a23e04dbc61dfd15d10c18d3b8668ad0ba19f3ba8d5e4afbf3968c51065ae192d5d3656134269ab7eaf71486a27e36b77e93f78e50732f3353199a36f639c6d7e1d268152921ddfb0a7f925a63acd0a250e0dc0e0222c367ce91732a02c4cf4875b3cccf701f2705762bd66d1cb0692df74e903b10979b7669bb9d1736eed0add399cb455f03b07682dc52a29e84782847f00af6b1896addaa5ad0ef6e2000e40b540200000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000054da8e1512b63608b1b4fa58485401057ca8ecedfa94970a64453683f84ec6548eda56cf0f769fdafab4233d19bd8976ae67826a5735e95069d069ca32daa15e209c5be447bae108000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000006400000000000000fd240300000002c4cf4875b3cccf701f2705762bd66d1cb0692df74e903b10979b7669bb9d173654da8e1512b63608b1b4fa58485401057ca8ecedfa94970a64453683f84ec6546a31e8caecb929d8f0b68a085bd5221c102a87eb8381d0a816cde5393274a83d2443161908961d62f49e39aba4082534343e089e812a885d30c19c4e8c28eb1fb26033c6d7ae4db3a36b21e9fb49327815669bf60109a55c62abd2fe72a39c62da9564ec818b9280bac512c861afea8ccc3453298fc596a8fcb08205e07df12a454f199887954a83a305a7f443960470a69151a8fc11e55bac7fa4ea02b39409e376f07423e1d521495f7b025e836ae2d691729fd7df710abe7dbf0d16f5db00a8cb360c18d125246e2120c937cd2407a1df3a41c7a1dd9723398469f7d9570bd7e4e15d1ee5f57d9c7d27fa11e7a87944988bb8729ce8a215df1b348846c20d5521d6ea013553b919fa0cc8160388a6b236c9a7a1c2b0d08f962b02bb004809025d0e6bd70b237210efa3026728d6d90300445d9658e921e343caafd747823ef28f6de349f2d593529886b1194a2c117ac8795bd9578b1ed217997017ce2832488a0df8a1c8b1a531c8c7f109297f0e2b4333cbdcb80f920bdb65de8a362e3f2cc1fa67665eee2aeddac1a3e665b0014238995ae8d7f513f1dd6e14cd115850e03b0fe9a7b5f86aa0a8d446de4d755a338c1883485607487a1f12941de1105fe22f41623a1fc76f44aeef58d20f4f42b8201b8eaca05b28aca3d2cf39f09836e460b63f5e9f2edda58da028879c3d933b3e24b57bf9c93d031be86c5b60be5ece9ca918ba6b181f15d37c668a25bec59a374540e761ea6a654bf989c89abf3c54cb920a9fb09cb56874664d93fa73e091ada2e321d1a469b324183f92dd66601cf3754ecf5a7a3d75ef8fcd702a20d7786de0bee1d43286b36d4e1c5af0be42142b0770d9409c7199cb1320fa965b4a45ba3ab20f24ec13d034432302e3da473c206ab0caf1d75ee5a86cf8dd4f3e5f2b897d0f504297521ea77e3319c9f77d74ab00fa79a3d9f9b17891c7e841b9bac9cbf14896f4d94456a243a7a0d8542d74768a420ff0a1bb15cf94a5b606d88e923741c45c8058f2deccbade1f22240590d0030000000000b914d8c4d31565c62f31b02add94b8fe33f13db5c6f7dec9e0fb363280e59801"
)

func DecodeGenesis() *block.Block {
	b := block.NewBlock()
	switch Get().General.Network {
	case "testnet":

		blob, err := hex.DecodeString(TestNetGenesisBlob)
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		buf.Write(blob)
		if err := block.Unmarshal(&buf, b); err != nil {
			panic(err)
		}
	}
	return b
}
