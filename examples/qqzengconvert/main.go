// asn-writer is an example of how to create an ASN MaxMind DB file from the
// GeoLite2 ASN CSVs. You must have the CSVs in the current working directory.
package main

import (
	"bufio"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"

	"./qqzeng"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

var continents_ch2en map[string]string
var continents_ch2code map[string]string
var province_ch2en map[string]string
var province_ch2code map[string]string
var city_ch2en map[string]string
var city_en2code map[string]string

func init() {
	continents_ch2en = make(map[string]string)
	continents_ch2en["亚洲"] = "Asia"
	continents_ch2en["北美洲"] = "North America"
	continents_ch2en["南极洲"] = "Antarctica"
	continents_ch2en["南美洲"] = "South America"
	continents_ch2en["大洋洲"] = "Oceania"
	continents_ch2en["欧洲"] = "Europe"
	continents_ch2en["非洲"] = "Africa"

	continents_ch2code = make(map[string]string)
	continents_ch2code["亚洲"] = "AS"
	continents_ch2code["北美洲"] = "NA"
	continents_ch2code["南极洲"] = "AN"
	continents_ch2code["南美洲"] = "SA"
	continents_ch2code["大洋洲"] = "OC"
	continents_ch2code["欧洲"] = "EU"
	continents_ch2code["非洲"] = "AF"

	province_ch2en = make(map[string]string)
	province_ch2en["YE:舍卜沃省"] = "Shabwah"
	province_ch2en["IR:礼萨呼罗珊省"] = "Razavi Khorasan"
	province_ch2en["IR:中央省"] = "Markazi"
	province_ch2en["AZ:纳希切万自治共和国"] = "Nakhichevan"
	province_ch2en["KE:曼德拉"] = "Mandera District"
	province_ch2en["KE:巴林戈郡"] = "Baringo"
	province_ch2en["CD:乔波省"] = "Tshopo"
	province_ch2en["CD:蒙加拉省"] = "Mongala"
	province_ch2en["CD:楚阿帕省"] = "Tshuapa"
	province_ch2en["DJ:塔朱拉州"] = "Tadjourah"
	province_ch2en["DJ:奧博克州"] = "Obock"
	province_ch2en["CF:上姆博穆省"] = "Haut-Mbomou"
	province_ch2en["CF:姆博穆省"] = "Mbomou"
	province_ch2en["CF:瓦卡加省"] = "Vakaga"
	province_ch2en["CF:上科托省"] = "Haute-Kotto"
	province_ch2en["CF:下科托省"] = "Basse-Kotto"
	province_ch2en["CF:瓦卡省"] = "Ouaka"
	province_ch2en["CF:巴明吉-班戈兰省"] = "Bamingui-Bangoran"
	province_ch2en["LB:南部省"] = "South Governorate"
	province_ch2en["BH:穆哈拉格省"] = "Muharraq"
	province_ch2en["BH:北方省"] = "Northern"
	province_ch2en["IL:海法区"] = "Haifa"
	province_ch2en["TR:马拉蒂亚省"] = "Malatya"
	province_ch2en["TR:穆拉省"] = "Muğla"
	province_ch2en["TR:阿达纳省"] = "Adana"
	province_ch2en["TR:约兹加特省"] = "Yozgat"
	province_ch2en["TR:哈塔伊省"] = "Hatay"
	province_ch2en["TR:安卡拉省"] = "Ankara"
	province_ch2en["TR:艾登省"] = "Aydın"
	province_ch2en["TR:伊斯帕尔塔省"] = "Isparta"
	province_ch2en["TR:凡城省"] = "Van"
	province_ch2en["TR:乌沙克省"] = "Uşak"
	province_ch2en["TR:马尼萨省"] = "Manisa"
	province_ch2en["TR:通杰利省"] = "Tunceli"
	province_ch2en["TR:屈塔希亚省"] = "Kütahya"
	province_ch2en["TR:比特利斯省"] = "Bitlis"
	province_ch2en["TR:梅尔辛省"] = "Mersin"
	province_ch2en["TR:开塞利省"] = "Kayseri"
	province_ch2en["TR:舍尔纳克省"] = "Şırnak"
	province_ch2en["TR:锡尔特省"] = "Siirt"
	province_ch2en["TR:安塔利亚省"] = "Antalya"
	province_ch2en["TR:代尼兹利省"] = "Denizli"
	province_ch2en["TR:马尔丁省"] = "Mardin"
	province_ch2en["TR:加济安泰普省"] = "Gaziantep"
	province_ch2en["TR:穆什省"] = "Muş"
	province_ch2en["TR:科尼亚省"] = "Konya"
	province_ch2en["TR:恰纳卡莱省"] = "Canakkale"
	province_ch2en["TR:克尔谢希尔省"] = "Kırşehir"
	province_ch2en["TR:克勒克卡莱省"] = "Kırıkkale"
	province_ch2en["TR:基利斯省"] = "Kilis"
	province_ch2en["TR:巴特曼省"] = "Batman"
	province_ch2en["TR:卡拉曼省"] = "Karaman"
	province_ch2en["TR:阿勒省"] = "Ağrı"
	province_ch2en["TR:卡赫拉曼马拉什省"] = "Kahramanmaraş"
	province_ch2en["TR:厄德尔省"] = "Iğdır"
	province_ch2en["TR:巴勒克埃西尔省"] = "Balıkesir"
	province_ch2en["TR:布尔杜尔省"] = "Burdur"
	province_ch2en["TR:埃斯基谢希尔省"] = "Eskişehir"
	province_ch2en["TR:埃尔祖鲁姆省"] = "Erzurum"
	province_ch2en["TR:埃尔津詹省"] = "Erzincan"
	province_ch2en["TR:哈卡里省"] = "Hakkâri"
	province_ch2en["TR:比莱吉克省"] = "Bilecik"
	province_ch2en["TR:宾格尔省"] = "Bingöl"
	province_ch2en["TR:阿克萨赖省"] = "Aksaray"
	province_ch2en["LV:文茨皮尔斯"] = "Ventspils"
	province_ch2en["LV:利耶帕亚"] = "Liepaja"
	province_ch2en["LV:叶尔加瓦"] = "Jelgava"
	province_ch2en["LV:陶格夫匹尔斯"] = "Daugavpils"
	province_ch2en["RU:莫斯科州"] = "Moscow Oblast"
	province_ch2en["RU:科米共和国"] = "Komi"
	province_ch2en["RU:乌里扬诺夫斯克州"] = "Ulyanovsk Oblast"
	province_ch2en["RU:摩爾曼斯克州"] = "Murmansk"
	province_ch2en["RU:达吉斯坦共和国"] = "Dagestan"
	province_ch2en["RU:楚瓦什共和国"] = "Chuvashia"
	province_ch2en["RU:卡累利阿共和国"] = "Karelia"
	province_ch2en["SE:北博滕"] = "Norrbotten County"
	province_ch2en["GE:卡赫季州"] = "Kakheti"
	province_ch2en["GE:伊梅列季州"] = "Imereti"
	province_ch2en["GE:姆茨赫塔-姆季阿涅季州"] = "Mtskheta-Mtianeti"
	province_ch2en["GE:克維莫-卡特利州"] = "Kvemo Kartli"
	province_ch2en["MD:布里切尼區"] = "Briceni"
	province_ch2en["FI:新地区"] = "Uusimaa"
	province_ch2en["UA:敖德萨州"] = "Odessa"
	province_ch2en["UA:基辅州"] = "Kyiv Oblast"
	province_ch2en["UA:波尔塔瓦州"] = "Poltava Oblast"
	province_ch2en["UA:赫尔松州"] = "Kherson Oblast"
	province_ch2en["HU:巴奇-基什孔州"] = "Bács-Kiskun"
	province_ch2en["TR:伊斯坦布尔"] = "Istanbul"
	province_ch2en["TR:科贾埃利省"] = "Kocaeli"
	province_ch2en["TR:亚洛瓦省"] = "Yalova"
	province_ch2en["TR:萨姆松省"] = "Samsun"
	province_ch2en["TR:吉雷松省"] = "Giresun"
	province_ch2en["TR:埃迪尔内省"] = "Edirne"
	province_ch2en["TR:锡诺普省"] = "Sinop"
	province_ch2en["TR:奥尔杜省"] = "Ordu"
	province_ch2en["TR:巴尔滕省"] = "Bartın"
	province_ch2en["TR:泰基尔达省"] = "Tekirdağ"
	province_ch2en["TR:特拉布宗省"] = "Trabzon"
	province_ch2en["TR:卡斯塔莫努省"] = "Kastamonu"
	province_ch2en["TR:博卢省"] = "Bolu"
	province_ch2en["TR:阿马西亚省"] = "Amasya"
	province_ch2en["TR:萨卡里亚省"] = "Sakarya"
	province_ch2en["TR:阿尔特温省"] = "Artvin"
	province_ch2en["TR:阿尔达汉省"] = "Ardahan"
	province_ch2en["TR:克尔克拉雷利省"] = "Kırklareli"
	province_ch2en["TR:乔鲁姆省"] = "Çorum"
	province_ch2en["TR:昌克勒省"] = "Çankırı"
	province_ch2en["PL:卢布林省"] = "Lublin"
	province_ch2en["PL:喀尔巴阡山省"] = "Subcarpathia"
	province_ch2en["PL:小波兰省"] = "Lesser Poland"
	province_ch2en["PL:罗兹省"] = "Łódź Voivodeship"
	province_ch2en["MK:耶古诺夫采区"] = "Jegunovce"
	province_ch2en["MK:博格丹奇区"] = "Bogdanci"
	province_ch2en["MK:雷森区"] = "Resen"
	province_ch2en["MK:多尔内尼区"] = "Dolneni"
	province_ch2en["MK:戈斯蒂瓦尔"] = "Gostivar"
	province_ch2en["MK:格拉德斯科区"] = "Gradsko"
	province_ch2en["AO:莫希科省"] = "Moxico"
	province_ch2en["KM:大科摩罗岛"] = "Grande Comore"
	province_ch2en["ZA:豪登省"] = "Gauteng"
	province_ch2en["ZA:林波波省"] = "Limpopo"
	province_ch2en["MZ:太特省"] = "Tete"
	province_ch2en["MZ:楠普拉省"] = "Nampula"
	province_ch2en["LS:塔巴-采卡區"] = "Thaba-Tseka"
	province_ch2en["AF:巴尔赫省"] = "Balkh"
	province_ch2en["PK:自由克什米爾"] = "Azad Jammu and Kashmir"
	province_ch2en["PK:信德省"] = "Sindh"
	province_ch2en["ID:北苏门答腊省"] = "North Sumatra"
	province_ch2en["TM:列巴普州"] = "Lebap"
	province_ch2en["MY:吉打"] = "Kedah"
	province_ch2en["LK:西北省"] = "North Western Province"
	province_ch2en["LK:中央省"] = "Central Province"
	province_ch2en["LK:北部省"] = "Northern Province"
	province_ch2en["IN:泰米尔纳德邦"] = "Tamil Nadu"
	province_ch2en["IN:拉贾斯坦邦"] = "Rajasthan"
	province_ch2en["IN:特伦甘纳邦"] = "Telangana"
	province_ch2en["IN:马哈拉施特拉邦"] = "Maharashtra"
	province_ch2en["IN:卡纳塔克邦"] = "Karnataka"
	province_ch2en["IN:安得拉邦"] = "Andhra Pradesh"
	province_ch2en["IN:哈里亚纳邦"] = "Haryana"
	province_ch2en["IN:西孟加拉邦"] = "West Bengal"
	province_ch2en["IN:古吉拉特邦"] = "Gujarat"
	province_ch2en["IN:北方邦"] = "Uttar Pradesh"
	province_ch2en["IN:中央邦"] = "Madhya Pradesh"
	province_ch2en["IN:果阿邦"] = "Goa"
	province_ch2en["IN:喀拉拉邦"] = "Kerala"
	province_ch2en["IN:喜马偕尔邦"] = "Himachal Pradesh"
	province_ch2en["IN:曼尼普尔邦"] = "Manipur"
	province_ch2en["IN:阿萨姆邦"] = "Assam"
	province_ch2en["IN:特里普拉邦"] = "Tripura"
	province_ch2en["IN:梅加拉亚邦"] = "Meghalaya"
	province_ch2en["IN:旁遮普邦"] = "Punjab"
	province_ch2en["IN:恰蒂斯加尔邦"] = "Chhattisgarh"
	province_ch2en["IN:比哈尔邦"] = "Bihar"
	province_ch2en["IN:锡金"] = "Sikkim"
	province_ch2en["IN:贾坎德邦"] = "Jharkhand"
	province_ch2en["IN:米佐拉姆邦"] = "Mizoram"
	province_ch2en["IN:阿鲁纳恰尔邦"] = "Arunachal Pradesh"
	province_ch2en["IN:安达曼-尼科巴群岛"] = "Andaman and Nicobar"
	province_ch2en["IN:本地治里"] = "Union Territory of Puducherry"
	province_ch2en["IN:那加兰邦"] = "Nagaland"
	province_ch2en["IN:拉達克"] = "Ladakh"
	province_ch2en["IN:昌迪加尔"] = "Chandigarh"
	province_ch2en["CN:云南"] = "Yunnan"
	province_ch2en["CN:甘肃"] = "Gansu"
	province_ch2en["CN:西藏自治区"] = "Tibet"
	province_ch2en["CN:新疆"] = "Xinjiang Uyghur Autonomous Region"
	province_ch2en["CN:江苏省"] = "Jiangsu"
	province_ch2en["CN:贵州"] = "Guizhou"
	province_ch2en["CN:安徽"] = "Anhui"
	province_ch2en["VN:胡志明市"] = "Ho Chi Minh"
	province_ch2en["ID:东爪哇省"] = "East Java"
	province_ch2en["ID:中爪哇省"] = "Central Java"
	province_ch2en["TL:維克克區"] = "Viqueque"
	province_ch2en["ID:巴厘岛"] = "Bali"
	province_ch2en["ID:北苏拉威西省"] = "North Sulawesi"
	province_ch2en["ID:邦加-勿里洞省"] = "Bangka–Belitung Islands"
	province_ch2en["ID:東加里曼丹省"] = "East Kalimantan"
	province_ch2en["ID:西爪哇省"] = "West Java"
	province_ch2en["ID:南加里曼丹省"] = "South Kalimantan"
	province_ch2en["ID:万丹省"] = "Banten"
	province_ch2en["ID:南苏门答腊省"] = "South Sumatra"
	province_ch2en["ID:西加里曼丹省"] = "West Kalimantan"
	province_ch2en["ID:中苏拉威西省"] = "Central Sulawesi"
	province_ch2en["TW:新北市"] = "New Taipei"
	province_ch2en["TW:苗栗縣"] = "Miaoli"
	province_ch2en["TW:台南市"] = "Tainan"
	province_ch2en["TW:高雄市"] = "Kaohsiung"
	province_ch2en["TW:台北市"] = "Taipei City"
	province_ch2en["TW:新竹市"] = "Hsinchu"
	province_ch2en["TW:基隆市"] = "Keelung"
	province_ch2en["PH:中央吕宋"] = "Central Luzon"
	province_ch2en["PH:西米沙鄢"] = "Western Visayas"
	province_ch2en["PH:中米沙鄢"] = "Central Visayas"
	province_ch2en["PH:北棉兰老"] = "Northern Mindanao"
	province_ch2en["PH:卡加延河谷"] = "Cagayan Valley"
	province_ch2en["PH:棉兰老穆斯林自治区"] = "Autonomous Region in Muslim Mindanao"
	province_ch2en["MY:柔佛州"] = "Johor"
	province_ch2en["MY:砂拉越"] = "Sarawak"
	province_ch2en["CN:山东省"] = "Shandong"
	province_ch2en["CN:陕西"] = "Shaanxi"
	province_ch2en["CN:四川省"] = "Sichuan"
	province_ch2en["CN:湖南"] = "Hunan"
	province_ch2en["CN:河北省"] = "Hebei"
	province_ch2en["CN:河南"] = "Henan"
	province_ch2en["CN:浙江省"] = "Zhejiang"
	province_ch2en["CN:辽宁"] = "Liaoning"
	province_ch2en["CN:宁夏回族自治区"] = "Ningxia Hui Autonomous Region"
	province_ch2en["CN:广东"] = "Guangdong"
	province_ch2en["CN:福建省"] = "Fujian"
	province_ch2en["CN:天津市"] = "Tianjin"
	province_ch2en["CN:上海"] = "Shanghai"
	province_ch2en["CN:湖北省"] = "Hubei"
	province_ch2en["CN:山西"] = "Shanxi"
	province_ch2en["CN:广西壮族自治区"] = "Guangxi"
	province_ch2en["CN:江西"] = "Jiangxi"
	province_ch2en["CN:青海省"] = "Qinghai"
	province_ch2en["CN:内蒙古自治区"] = "Inner Mongolia Autonomous Region"
	province_ch2en["CN:北京市"] = "Beijing"
	province_ch2en["CN:海南"] = "Hainan"
	province_ch2en["CN:重庆"] = "Chongqing"
	province_ch2en["KH:班迭棉吉省"] = "Banteay Meanchey"
	province_ch2en["KH:菩萨省"] = "Pursat"
	province_ch2en["KH:干丹省"] = "Kandal"
	province_ch2en["KH:柏威夏省"] = "Preah Vihear"
	province_ch2en["KH:戈公省"] = "Koh Kong"
	province_ch2en["KH:贡布省"] = "Kampot"
	province_ch2en["KH:西哈努克市"] = "Preah Sihanouk"
	province_ch2en["KH:磅清扬省"] = "Kampong Chhnang"
	province_ch2en["KH:磅湛省"] = "Kampong Cham"
	province_ch2en["KH:腊塔纳基里省"] = "Ratanakiri"
	province_ch2en["KR:首尔特别市"] = "Seoul"
	province_ch2en["KR:京畿道"] = "Gyeonggi-do"
	province_ch2en["KR:釜山广域市"] = "Busan"
	province_ch2en["KR:光州广域市"] = "Gwangju"
	province_ch2en["JP:和歌山县"] = "Wakayama"
	province_ch2en["JP:兵库县"] = "Hyōgo"
	province_ch2en["JP:新潟县"] = "Niigata"
	province_ch2en["JP:山形县"] = "Yamagata"
	province_ch2en["JP:岐阜县"] = "Gifu"
	province_ch2en["JP:福冈县"] = "Fukuoka"
	province_ch2en["JP:山梨县"] = "Yamanashi"
	province_ch2en["JP:滋贺县"] = "Shiga"
	province_ch2en["JP:岛根县"] = "Shimane"
	province_ch2en["JP:山口县"] = "Yamaguchi"
	province_ch2en["JP:福岛县"] = "Fukushima-ken"
	province_ch2en["JP:长野县"] = "Nagano"
	province_ch2en["JP:香川县"] = "Kagawa"
	province_ch2en["JP:秋田县"] = "Akita"
	province_ch2en["RU:犹太自治州"] = "Yevrey (Jewish) Autonomous Oblast"
	province_ch2en["CN:吉林"] = "Jilin"
	province_ch2en["CN:黑龙江省"] = "Heilongjiang"
	province_ch2en["JP:北海道"] = "Hokkaido"
	province_ch2en["JP:岩手县"] = "Iwate"
	province_ch2en["AU:澳大利亚首都领地"] = "Australian Capital Territory"
	province_ch2en["NZ:南地大区"] = "Southland"
	province_ch2en["NZ:北地大区"] = "Northland"
	province_ch2en["NZ:惠灵顿"] = "Wellington"
	province_ch2en["NZ:西岸大区"] = "West Coast"
	province_ch2en["NZ:奥塔哥大区"] = "Otago"
	province_ch2en["NZ:马尔堡"] = "Marlborough"
	province_ch2en["NZ:普伦蒂湾大区"] = "Bay of Plenty"
	province_ch2en["FJ:中央大区"] = "Central"
	province_ch2en["FJ:西部大区"] = "Western"
	province_ch2en["FJ:北部大区"] = "Northern"
	province_ch2en["NZ:吉斯伯恩大区"] = "Gisborne"
	province_ch2en["LY:朱夫拉省"] = "Al Jufrah"
	province_ch2en["SN:济金绍尔区"] = "Ziguinchor"
	province_ch2en["SN:久尔贝勒区"] = "Diourbel"
	province_ch2en["SN:达喀尔区"] = "Dakar"
	province_ch2en["SN:考拉克区"] = "Kaolack"
	province_ch2en["CG:盆地省"] = "Cuvette"
	province_ch2en["CG:利夸拉省"] = "Likouala"
	province_ch2en["PT:里斯本區"] = "Lisbon"
	province_ch2en["CD:奎卢省"] = "Kwilu"
	province_ch2en["TG:中部区"] = "Centrale"
	province_ch2en["TG:卡拉区"] = "Kara"
	province_ch2en["MR:塔甘特省"] = "Tagant"
	province_ch2en["MR:提里斯-宰穆爾省"] = "Tiris Zemmour"
	province_ch2en["MR:卜拉克納省"] = "Brakna"
	province_ch2en["MR:因希里省"] = "Inchiri"
	province_ch2en["CF:洛巴耶省"] = "Lobaye"
	province_ch2en["CF:瓦姆省"] = "Ouham"
	province_ch2en["BJ:阿黎博里省"] = "Alibori"
	province_ch2en["BJ:峽谷省"] = "Donga"
	province_ch2en["GA:河口省"] = "Estuaire"
	province_ch2en["ST:普林西比岛"] = "Principe"
	province_ch2en["NE:迪法大区"] = "Diffa"
	province_ch2en["ES:埃斯特雷马杜拉"] = "Extremadura"
	province_ch2en["ES:卡斯蒂利亚-拉曼恰"] = "Castille-La Mancha"
	province_ch2en["ES:加那利群岛"] = "Canary Islands"
	province_ch2en["IT:西西里岛"] = "Sicily"
	province_ch2en["IT:卡拉布里亚"] = "Calabria"
	province_ch2en["IT:巴斯利卡塔"] = "Basilicate"
	province_ch2en["IT:普利亚"] = "Apulia"
	province_ch2en["IT:坎帕尼亚"] = "Campania"
	province_ch2en["DK:西兰大区"] = "Zealand"
	province_ch2en["DK:南丹麦大区"] = "South Denmark"
	province_ch2en["DK:中日德兰大区"] = "Central Jutland"
	province_ch2en["GB:英格兰"] = "England"
	province_ch2en["GB:苏格兰"] = "Scotland"
	province_ch2en["CH:汝拉州"] = "Jura"
	province_ch2en["SE:厄勒布鲁省"] = "Örebro County"
	province_ch2en["SE:西约塔兰省"] = "Västra Götaland County"
	province_ch2en["SE:耶夫勒堡省"] = "Gävleborg County"
	province_ch2en["SJ:斯瓦尔巴群岛"] = "Svalbard"
	province_ch2en["NL:北布拉班特省"] = "North Brabant"
	province_ch2en["AT:克恩顿州"] = "Carinthia"
	province_ch2en["AT:福拉尔贝格州"] = "Vorarlberg"
	province_ch2en["AT:布尔根兰州"] = "Burgenland"
	province_ch2en["DE:巴伐利亚"] = "Bavaria"
	province_ch2en["DE:巴登-符腾堡"] = "Baden-Württemberg"
	province_ch2en["DE:莱茵兰-普法尔茨"] = "Rheinland-Pfalz"
	province_ch2en["DE:勃兰登堡"] = "Brandenburg"
	province_ch2en["DE:萨克森-安哈尔特"] = "Saxony-Anhalt"
	province_ch2en["DE:石勒苏益格-荷尔斯泰因"] = "Schleswig-Holstein"
	province_ch2en["DE:下萨克森"] = "Lower Saxony"
	province_ch2en["DE:柏林"] = "Land Berlin"
	province_ch2en["IE:芒斯特省"] = "Munster"
	province_ch2en["FR:卢瓦尔河地区"] = "Pays de la Loire"
	province_ch2en["FR:布列塔尼半岛"] = "Brittany"
	province_ch2en["FR:普罗旺斯-阿尔卑斯-蓝色海岸"] = "Provence-Alpes-Côte d'Azur"
	province_ch2en["AD:卡尼略"] = "Canillo"
	province_ch2en["AD:圣胡利娅-德洛里亚"] = "Sant Julià de Loria"
	province_ch2en["AD:奥尔迪诺"] = "Ordino"
	province_ch2en["LI:毛伦"] = "Mauren"
	province_ch2en["LI:埃申"] = "Eschen"
	province_ch2en["HU:佐洛州"] = "Zala"
	province_ch2en["HU:沃什州"] = "Vas"
	province_ch2en["CZ:中波希米亚州"] = "Central Bohemia"
	province_ch2en["CZ:南波希米亚州"] = "Jihocesky kraj"
	province_ch2en["PL:西里西亚省"] = "Silesia"
	province_ch2en["PL:大波兰省"] = "Greater Poland"
	province_ch2en["PL:下西里西亚省"] = "Lower Silesia"
	province_ch2en["ES:卡斯蒂利亚-莱昂"] = "Castille and León"
	province_ch2en["ES:加利西亚"] = "Galicia"
	province_ch2en["ES:坎塔布里亚"] = "Cantabria"
	province_ch2en["ES:拉里奥哈"] = "La Rioja"
	province_ch2en["IT:威尼托"] = "Veneto"
	province_ch2en["IT:利古里亞"] = "Liguria"
	province_ch2en["IT:拉齐奥"] = "Latium"
	province_ch2en["IT:皮埃蒙特"] = "Piedmont"
	province_ch2en["IT:托斯卡纳"] = "Tuscany"
	province_ch2en["IT:阿布鲁佐"] = "Abruzzo"
	province_ch2en["IT:马尔凯"] = "The Marches"
	province_ch2en["IT:莫利塞"] = "Molise"
	province_ch2en["IT:翁布里亚"] = "Umbria"
	province_ch2en["SM:基埃萨努欧瓦"] = "Chiesanuova"
	province_ch2en["SI:伊德里亞"] = "Idrija"
	province_ch2en["SI:斯洛文尼亞科尼采"] = "Slovenske Konjice"
	province_ch2en["SI:斯洛文尼亞比斯特里察"] = "Slovenska Bistrica"
	province_ch2en["SI:伊利爾斯卡比斯特里察"] = "Ilirska Bistrica"
	province_ch2en["SI:拉多夫利察"] = "Radovljica"
	province_ch2en["SI:科佩尔"] = "Koper"
	province_ch2en["SI:皮夫卡"] = "Pivka"
	province_ch2en["SI:普图伊"] = "Ptuj"
	province_ch2en["SI:普雷瓦列"] = "Prevalje"
	province_ch2en["ME:莫伊科瓦茨"] = "Mojkovac"
	province_ch2en["SI:马里博尔"] = "Maribor"
	province_ch2en["SI:伊佐拉"] = "Izola"
	province_ch2en["ME:新海尔采格"] = "Herceg Novi"
	province_ch2en["SI:赫拉斯特尼克"] = "Hrastnik"
	province_ch2en["SI:洛加泰茨"] = "Logatec"
	province_ch2en["ME:比耶洛波列"] = "Bijelo Polje"
	province_ch2en["SI:采爾克諾"] = "Cerkno"
	province_ch2en["AO:万博省"] = "Huambo"
	province_ch2en["AO:本吉拉省"] = "Benguela"
	province_ch2en["BR:伯南布哥州"] = "Pernambuco"
	province_ch2en["BR:托坎廷斯州"] = "Tocantins"
	province_ch2en["BR:马拉尼昂州"] = "Maranhao"
	province_ch2en["BR:帕拉州"] = "Para"
	province_ch2en["BR:亚马孙州"] = "Amazonas"
	province_ch2en["GL:凯克卡塔"] = "Qeqqata"
	province_ch2en["GL:库雅雷克"] = "Kujalleq"
	province_ch2en["AR:布宜诺斯艾利斯"] = "Buenos Aires F.D."
	province_ch2en["UY:塞罗拉尔戈省"] = "Cerro Largo"
	province_ch2en["UY:卡内洛内斯省"] = "Canelones"
	province_ch2en["UY:阿蒂加斯省"] = "Artigas"
	province_ch2en["BR:米纳斯吉拉斯州"] = "Minas Gerais"
	province_ch2en["BR:圣埃斯皮里图州"] = "Espirito Santo"
	province_ch2en["BR:戈亚斯州"] = "Goias"
	province_ch2en["BR:马托格罗索州"] = "Mato Grosso"
	province_ch2en["BR:南马托格罗索州"] = "Mato Grosso do Sul"
	province_ch2en["MX:塔毛利帕斯州"] = "Tamaulipas"
	province_ch2en["MX:墨西哥州"] = "México"
	province_ch2en["MX:瓦哈卡州"] = "Oaxaca"
	province_ch2en["MX:普埃布拉州"] = "Puebla"
	province_ch2en["MX:莫雷洛斯州"] = "Morelos"
	province_ch2en["MX:尤卡坦州"] = "Yucatán"
	province_ch2en["MX:恰帕斯州"] = "Chiapas"
	province_ch2en["MX:墨西哥城市"] = "Mexico City"
	province_ch2en["MX:塔巴斯科州"] = "Tabasco"
	province_ch2en["MX:克雷塔羅州"] = "Querétaro"
	province_ch2en["MX:坎佩切州"] = "Campeche"
	province_ch2en["VE:阿马库罗三角洲州"] = "Delta Amacuro"
	province_ch2en["PE:安卡什大区"] = "Ancash"
	province_ch2en["PE:洛雷托大区"] = "Loreto"
	province_ch2en["PE:拉利伯塔德大区"] = "La Libertad"
	province_ch2en["PE:乌卡亚利大区"] = "Ucayali"
	province_ch2en["PE:通贝斯大区"] = "Tumbes"
	province_ch2en["PE:兰巴耶克大区"] = "Lambayeque"
	province_ch2en["PE:皮乌拉地区"] = "Piura"
	province_ch2en["PE:亚马孙大区"] = "Amazonas"
	province_ch2en["PA:恩贝拉-沃内安特区"] = "Embera-Wounaan"
	province_ch2en["PA:库纳雅拉特区"] = "Guna Yala"
	province_ch2en["CL:圣地亚哥首都大区"] = "Santiago Metropolitan"
	province_ch2en["CL:塔拉帕卡大区"] = "Tarapacá"
	province_ch2en["PE:帕斯科大区"] = "Pasco"
	province_ch2en["PE:阿雷基帕大区"] = "Arequipa"
	province_ch2en["PE:胡宁大区"] = "Junin"
	province_ch2en["PE:马德雷德迪奥斯大区"] = "Madre de Dios"
	province_ch2en["PE:塔克纳大区"] = "Tacna"
	province_ch2en["PE:阿亚库乔大区"] = "Ayacucho"
	province_ch2en["PE:伊卡大区"] = "Ica"
	province_ch2en["PE:普诺大区"] = "Puno"
	province_ch2en["MX:奇瓦瓦州"] = "Chihuahua"
	province_ch2en["MX:哈利斯科州"] = "Jalisco"
	province_ch2en["MX:米却肯州"] = "Michoacán"
	province_ch2en["WS:艾加伊勒泰"] = "Aiga-i-le-Tai"
	province_ch2en["WS:阿图阿"] = "Atua"
	province_ch2en["WS:图阿马萨加"] = "Tuamasaga"
	province_ch2en["FJ:东部大区"] = "Eastern"
	province_ch2en["US:德克萨斯州"] = "Texas"
	province_ch2en["US:亚拉巴马州"] = "Alabama"
	province_ch2en["US:弗吉尼亚州"] = "Virginia"
	province_ch2en["US:阿肯色州"] = "Arkansas"
	province_ch2en["US:特拉华州"] = "Delaware"
	province_ch2en["US:佛罗里达州"] = "Florida"
	province_ch2en["US:乔治亚"] = "Georgia"
	province_ch2en["US:伊利诺伊州"] = "Illinois"
	province_ch2en["US:印第安纳州"] = "Indiana"
	province_ch2en["US:马里兰州"] = "Maryland"
	province_ch2en["US:肯塔基州"] = "Kentucky"
	province_ch2en["US:密苏里州"] = "Missouri"
	province_ch2en["US:密西西比州"] = "Mississippi"
	province_ch2en["US:北卡罗来纳州"] = "North Carolina"
	province_ch2en["US:南卡罗来纳州"] = "South Carolina"
	province_ch2en["US:田纳西州"] = "Tennessee"
	province_ch2en["US:路易斯安那州"] = "Louisiana"
	province_ch2en["US:新泽西州"] = "New Jersey"
	province_ch2en["US:俄亥俄州"] = "Ohio"
	province_ch2en["US:宾夕法尼亚州"] = "Pennsylvania"
	province_ch2en["US:康乃狄克州"] = "Connecticut"
	province_ch2en["US:艾奥瓦州"] = "Iowa"
	province_ch2en["US:缅因州"] = "Maine"
	province_ch2en["US:密歇根州"] = "Michigan"
	province_ch2en["US:纽约州"] = "New York"
	province_ch2en["US:南达科他州"] = "South Dakota"
	province_ch2en["US:威斯康辛州"] = "Wisconsin"
	province_ch2en["US:明尼苏达州"] = "Minnesota"
	province_ch2en["US:北达科他州"] = "North Dakota"
	province_ch2en["US:新罕布什尔州"] = "New Hampshire"
	province_ch2en["US:佛蒙特州"] = "Vermont"
	province_ch2en["US:加利福尼亚州"] = "California"
	province_ch2en["US:新墨西哥州"] = "New Mexico"
	province_ch2en["US:犹他州"] = "Utah"
	province_ch2en["US:内华达州"] = "Nevada"
	province_ch2en["US:爱达荷州"] = "Idaho"
	province_ch2en["US:阿拉斯加州"] = "Alaska"
	province_ch2en["US:蒙大拿州"] = "Montana"
	province_ch2en["US:俄勒冈州"] = "Oregon"
	province_ch2en["US:华盛顿州"] = "Washington"
	province_ch2en["US:怀俄明州"] = "Wyoming"
	province_ch2en["US:夏威夷州"] = "Hawaii"
	province_ch2en["CA:不列颠哥伦比亚"] = "British Columbia"
	province_ch2en["CA:艾伯塔"] = "Alberta"
	province_ch2en["CA:安大略"] = "Ontario"
	province_ch2en["CA:新斯科舍"] = "Nova Scotia"
	province_ch2en["CA:曼尼托巴"] = "Manitoba"
	province_ch2en["CA:新不倫瑞克"] = "New Brunswick"
	province_ch2en["CA:育空"] = "Yukon"
	province_ch2en["ES:休达"] = "Ceuta"
	province_ch2en["MY:布城"] = "Putrajaya"

	province_ch2code = make(map[string]string)
	province_ch2code["country_iso_code:subdivision_1_name"] = "subdivision_1_iso_code"
	province_ch2code["YE:舍卜沃省"] = "SH"
	province_ch2code["IR:礼萨呼罗珊省"] = "09"
	province_ch2code["IR:中央省"] = "00"
	province_ch2code["AZ:纳希切万自治共和国"] = "NX"
	province_ch2code["KE:曼德拉"] = "24"
	province_ch2code["KE:巴林戈郡"] = "01"
	province_ch2code["CD:乔波省"] = "TO"
	province_ch2code["CD:蒙加拉省"] = "MO"
	province_ch2code["CD:楚阿帕省"] = "TU"
	province_ch2code["DJ:塔朱拉州"] = "TA"
	province_ch2code["DJ:奧博克州"] = "OB"
	province_ch2code["CF:上姆博穆省"] = "HM"
	province_ch2code["CF:姆博穆省"] = "MB"
	province_ch2code["CF:瓦卡加省"] = "VK"
	province_ch2code["CF:上科托省"] = "HK"
	province_ch2code["CF:下科托省"] = "BK"
	province_ch2code["CF:瓦卡省"] = "UK"
	province_ch2code["CF:巴明吉-班戈兰省"] = "BB"
	province_ch2code["LB:南部省"] = "JA"
	province_ch2code["BH:穆哈拉格省"] = "15"
	province_ch2code["BH:北方省"] = "17"
	province_ch2code["IL:海法区"] = "HA"
	province_ch2code["TR:马拉蒂亚省"] = "44"
	province_ch2code["TR:穆拉省"] = "48"
	province_ch2code["TR:阿达纳省"] = "01"
	province_ch2code["TR:约兹加特省"] = "66"
	province_ch2code["TR:哈塔伊省"] = "31"
	province_ch2code["TR:安卡拉省"] = "06"
	province_ch2code["TR:艾登省"] = "09"
	province_ch2code["TR:伊斯帕尔塔省"] = "32"
	province_ch2code["TR:凡城省"] = "65"
	province_ch2code["TR:乌沙克省"] = "64"
	province_ch2code["TR:马尼萨省"] = "45"
	province_ch2code["TR:通杰利省"] = "62"
	province_ch2code["TR:屈塔希亚省"] = "43"
	province_ch2code["TR:比特利斯省"] = "13"
	province_ch2code["TR:梅尔辛省"] = "33"
	province_ch2code["TR:开塞利省"] = "38"
	province_ch2code["TR:舍尔纳克省"] = "73"
	province_ch2code["TR:锡尔特省"] = "56"
	province_ch2code["TR:安塔利亚省"] = "07"
	province_ch2code["TR:代尼兹利省"] = "20"
	province_ch2code["TR:马尔丁省"] = "47"
	province_ch2code["TR:加济安泰普省"] = "27"
	province_ch2code["TR:穆什省"] = "49"
	province_ch2code["TR:科尼亚省"] = "42"
	province_ch2code["TR:恰纳卡莱省"] = "17"
	province_ch2code["TR:克尔谢希尔省"] = "40"
	province_ch2code["TR:克勒克卡莱省"] = "71"
	province_ch2code["TR:基利斯省"] = "79"
	province_ch2code["TR:巴特曼省"] = "72"
	province_ch2code["TR:卡拉曼省"] = "70"
	province_ch2code["TR:阿勒省"] = "04"
	province_ch2code["TR:卡赫拉曼马拉什省"] = "46"
	province_ch2code["TR:厄德尔省"] = "76"
	province_ch2code["TR:巴勒克埃西尔省"] = "10"
	province_ch2code["TR:布尔杜尔省"] = "15"
	province_ch2code["TR:埃斯基谢希尔省"] = "26"
	province_ch2code["TR:埃尔祖鲁姆省"] = "25"
	province_ch2code["TR:埃尔津詹省"] = "24"
	province_ch2code["TR:哈卡里省"] = "30"
	province_ch2code["TR:比莱吉克省"] = "11"
	province_ch2code["TR:宾格尔省"] = "12"
	province_ch2code["TR:阿克萨赖省"] = "68"
	province_ch2code["LV:文茨皮尔斯"] = "VEN"
	province_ch2code["LV:利耶帕亚"] = "LPX"
	province_ch2code["LV:叶尔加瓦"] = "JEL"
	province_ch2code["LV:陶格夫匹尔斯"] = "DGV"
	province_ch2code["RU:莫斯科州"] = "MOS"
	province_ch2code["RU:科米共和国"] = "KO"
	province_ch2code["RU:乌里扬诺夫斯克州"] = "ULY"
	province_ch2code["RU:摩爾曼斯克州"] = "MUR"
	province_ch2code["RU:达吉斯坦共和国"] = "DA"
	province_ch2code["RU:楚瓦什共和国"] = "CU"
	province_ch2code["RU:卡累利阿共和国"] = "KR"
	province_ch2code["SE:北博滕"] = "BD"
	province_ch2code["GE:卡赫季州"] = "KA"
	province_ch2code["GE:伊梅列季州"] = "IM"
	province_ch2code["GE:姆茨赫塔-姆季阿涅季州"] = "MM"
	province_ch2code["GE:克維莫-卡特利州"] = "KK"
	province_ch2code["MD:布里切尼區"] = "BR"
	province_ch2code["FI:新地区"] = "18"
	province_ch2code["UA:敖德萨州"] = "51"
	province_ch2code["UA:基辅州"] = "32"
	province_ch2code["UA:波尔塔瓦州"] = "53"
	province_ch2code["UA:赫尔松州"] = "65"
	province_ch2code["HU:巴奇-基什孔州"] = "BK"
	province_ch2code["TR:伊斯坦布尔"] = "34"
	province_ch2code["TR:科贾埃利省"] = "41"
	province_ch2code["TR:亚洛瓦省"] = "77"
	province_ch2code["TR:萨姆松省"] = "55"
	province_ch2code["TR:吉雷松省"] = "28"
	province_ch2code["TR:埃迪尔内省"] = "22"
	province_ch2code["TR:锡诺普省"] = "57"
	province_ch2code["TR:奥尔杜省"] = "52"
	province_ch2code["TR:巴尔滕省"] = "74"
	province_ch2code["TR:泰基尔达省"] = "59"
	province_ch2code["TR:特拉布宗省"] = "61"
	province_ch2code["TR:卡斯塔莫努省"] = "37"
	province_ch2code["TR:博卢省"] = "14"
	province_ch2code["TR:阿马西亚省"] = "05"
	province_ch2code["TR:萨卡里亚省"] = "54"
	province_ch2code["TR:阿尔特温省"] = "08"
	province_ch2code["TR:阿尔达汉省"] = "75"
	province_ch2code["TR:克尔克拉雷利省"] = "39"
	province_ch2code["TR:乔鲁姆省"] = "19"
	province_ch2code["TR:昌克勒省"] = "18"
	province_ch2code["PL:卢布林省"] = "06"
	province_ch2code["PL:喀尔巴阡山省"] = "18"
	province_ch2code["PL:小波兰省"] = "12"
	province_ch2code["PL:罗兹省"] = "10"
	province_ch2code["MK:耶古诺夫采区"] = "606"
	province_ch2code["MK:博格丹奇区"] = "401"
	province_ch2code["MK:雷森区"] = "509"
	province_ch2code["MK:多尔内尼区"] = "503"
	province_ch2code["MK:戈斯蒂瓦尔"] = "604"
	province_ch2code["MK:格拉德斯科区"] = "102"
	province_ch2code["AO:莫希科省"] = "MOX"
	province_ch2code["KM:大科摩罗岛"] = "G"
	province_ch2code["ZA:豪登省"] = "GP"
	province_ch2code["ZA:林波波省"] = "LP"
	province_ch2code["MZ:太特省"] = "T"
	province_ch2code["MZ:楠普拉省"] = "N"
	province_ch2code["LS:塔巴-采卡區"] = "K"
	province_ch2code["AF:巴尔赫省"] = "BAL"
	province_ch2code["PK:自由克什米爾"] = "JK"
	province_ch2code["PK:信德省"] = "SD"
	province_ch2code["ID:北苏门答腊省"] = "SU"
	province_ch2code["TM:列巴普州"] = "L"
	province_ch2code["MY:吉打"] = "02"
	province_ch2code["LK:西北省"] = "6"
	province_ch2code["LK:中央省"] = "2"
	province_ch2code["LK:北部省"] = "4"
	province_ch2code["IN:泰米尔纳德邦"] = "TN"
	province_ch2code["IN:拉贾斯坦邦"] = "RJ"
	province_ch2code["IN:特伦甘纳邦"] = "TG"
	province_ch2code["IN:马哈拉施特拉邦"] = "MH"
	province_ch2code["IN:卡纳塔克邦"] = "KA"
	province_ch2code["IN:安得拉邦"] = "AP"
	province_ch2code["IN:哈里亚纳邦"] = "HR"
	province_ch2code["IN:西孟加拉邦"] = "WB"
	province_ch2code["IN:古吉拉特邦"] = "GJ"
	province_ch2code["IN:北方邦"] = "UP"
	province_ch2code["IN:中央邦"] = "MP"
	province_ch2code["IN:果阿邦"] = "GA"
	province_ch2code["IN:喀拉拉邦"] = "KL"
	province_ch2code["IN:喜马偕尔邦"] = "HP"
	province_ch2code["IN:曼尼普尔邦"] = "MN"
	province_ch2code["IN:阿萨姆邦"] = "AS"
	province_ch2code["IN:特里普拉邦"] = "TR"
	province_ch2code["IN:梅加拉亚邦"] = "ML"
	province_ch2code["IN:旁遮普邦"] = "PB"
	province_ch2code["IN:恰蒂斯加尔邦"] = "CT"
	province_ch2code["IN:比哈尔邦"] = "BR"
	province_ch2code["IN:锡金"] = "SK"
	province_ch2code["IN:贾坎德邦"] = "JH"
	province_ch2code["IN:米佐拉姆邦"] = "MZ"
	province_ch2code["IN:阿鲁纳恰尔邦"] = "AR"
	province_ch2code["IN:安达曼-尼科巴群岛"] = "AN"
	province_ch2code["IN:本地治里"] = "PY"
	province_ch2code["IN:那加兰邦"] = "NL"
	province_ch2code["IN:拉達克"] = "LA"
	province_ch2code["IN:昌迪加尔"] = "CH"
	province_ch2code["CN:云南"] = "YN"
	province_ch2code["CN:甘肃"] = "GS"
	province_ch2code["CN:西藏自治区"] = "XZ"
	province_ch2code["CN:新疆"] = "XJ"
	province_ch2code["CN:江苏省"] = "JS"
	province_ch2code["CN:贵州"] = "GZ"
	province_ch2code["CN:安徽"] = "AH"
	province_ch2code["VN:胡志明市"] = "SG"
	province_ch2code["ID:东爪哇省"] = "JI"
	province_ch2code["ID:中爪哇省"] = "JT"
	province_ch2code["TL:維克克區"] = "VI"
	province_ch2code["ID:巴厘岛"] = "BA"
	province_ch2code["ID:北苏拉威西省"] = "SA"
	province_ch2code["ID:邦加-勿里洞省"] = "BB"
	province_ch2code["ID:東加里曼丹省"] = "KI"
	province_ch2code["ID:西爪哇省"] = "JB"
	province_ch2code["ID:南加里曼丹省"] = "KS"
	province_ch2code["ID:万丹省"] = "BT"
	province_ch2code["ID:南苏门答腊省"] = "SS"
	province_ch2code["ID:西加里曼丹省"] = "KB"
	province_ch2code["ID:中苏拉威西省"] = "ST"
	province_ch2code["TW:新北市"] = "NWT"
	province_ch2code["TW:苗栗縣"] = "MIA"
	province_ch2code["TW:台南市"] = "TNN"
	province_ch2code["TW:高雄市"] = "KHH"
	province_ch2code["TW:台北市"] = "TPE"
	province_ch2code["TW:新竹市"] = "HSQ"
	province_ch2code["TW:基隆市"] = "KEE"
	province_ch2code["PH:中央吕宋"] = "03"
	province_ch2code["PH:西米沙鄢"] = "06"
	province_ch2code["PH:中米沙鄢"] = "07"
	province_ch2code["PH:北棉兰老"] = "10"
	province_ch2code["PH:卡加延河谷"] = "02"
	province_ch2code["PH:棉兰老穆斯林自治区"] = "14"
	province_ch2code["MY:柔佛州"] = "01"
	province_ch2code["MY:砂拉越"] = "13"
	province_ch2code["CN:山东省"] = "SD"
	province_ch2code["CN:陕西"] = "SN"
	province_ch2code["CN:四川省"] = "SC"
	province_ch2code["CN:湖南"] = "HN"
	province_ch2code["CN:河北省"] = "HE"
	province_ch2code["CN:河南"] = "HA"
	province_ch2code["CN:浙江省"] = "ZJ"
	province_ch2code["CN:辽宁"] = "LN"
	province_ch2code["CN:宁夏回族自治区"] = "NX"
	province_ch2code["CN:广东"] = "GD"
	province_ch2code["CN:福建省"] = "FJ"
	province_ch2code["CN:天津市"] = "TJ"
	province_ch2code["CN:上海"] = "SH"
	province_ch2code["CN:湖北省"] = "HB"
	province_ch2code["CN:山西"] = "SX"
	province_ch2code["CN:广西壮族自治区"] = "GX"
	province_ch2code["CN:江西"] = "JX"
	province_ch2code["CN:青海省"] = "QH"
	province_ch2code["CN:内蒙古自治区"] = "NM"
	province_ch2code["CN:北京市"] = "BJ"
	province_ch2code["CN:海南"] = "HI"
	province_ch2code["CN:重庆"] = "CQ"
	province_ch2code["KH:班迭棉吉省"] = "1"
	province_ch2code["KH:菩萨省"] = "15"
	province_ch2code["KH:干丹省"] = "8"
	province_ch2code["KH:柏威夏省"] = "13"
	province_ch2code["KH:戈公省"] = "9"
	province_ch2code["KH:贡布省"] = "7"
	province_ch2code["KH:西哈努克市"] = "18"
	province_ch2code["KH:磅清扬省"] = "4"
	province_ch2code["KH:磅湛省"] = "3"
	province_ch2code["KH:腊塔纳基里省"] = "16"
	province_ch2code["KR:首尔特别市"] = "11"
	province_ch2code["KR:京畿道"] = "41"
	province_ch2code["KR:釜山广域市"] = "26"
	province_ch2code["KR:光州广域市"] = "29"
	province_ch2code["JP:和歌山县"] = "30"
	province_ch2code["JP:兵库县"] = "28"
	province_ch2code["JP:新潟县"] = "15"
	province_ch2code["JP:山形县"] = "06"
	province_ch2code["JP:岐阜县"] = "21"
	province_ch2code["JP:福冈县"] = "40"
	province_ch2code["JP:山梨县"] = "19"
	province_ch2code["JP:滋贺县"] = "25"
	province_ch2code["JP:岛根县"] = "32"
	province_ch2code["JP:山口县"] = "35"
	province_ch2code["JP:福岛县"] = "07"
	province_ch2code["JP:长野县"] = "20"
	province_ch2code["JP:香川县"] = "37"
	province_ch2code["JP:秋田县"] = "05"
	province_ch2code["RU:犹太自治州"] = "YEV"
	province_ch2code["CN:吉林"] = "JL"
	province_ch2code["CN:黑龙江省"] = "HL"
	province_ch2code["JP:北海道"] = "01"
	province_ch2code["JP:岩手县"] = "03"
	province_ch2code["AU:澳大利亚首都领地"] = "ACT"
	province_ch2code["NZ:南地大区"] = "STL"
	province_ch2code["NZ:北地大区"] = "NTL"
	province_ch2code["NZ:惠灵顿"] = "WGN"
	province_ch2code["NZ:西岸大区"] = "WTC"
	province_ch2code["NZ:奥塔哥大区"] = "OTA"
	province_ch2code["NZ:马尔堡"] = "MBH"
	province_ch2code["NZ:普伦蒂湾大区"] = "BOP"
	province_ch2code["FJ:中央大区"] = "C"
	province_ch2code["FJ:西部大区"] = "W"
	province_ch2code["FJ:北部大区"] = "N"
	province_ch2code["NZ:吉斯伯恩大区"] = "GIS"
	province_ch2code["LY:朱夫拉省"] = "JU"
	province_ch2code["SN:济金绍尔区"] = "ZG"
	province_ch2code["SN:久尔贝勒区"] = "DB"
	province_ch2code["SN:达喀尔区"] = "DK"
	province_ch2code["SN:考拉克区"] = "KL"
	province_ch2code["CG:盆地省"] = "8"
	province_ch2code["CG:利夸拉省"] = "7"
	province_ch2code["PT:里斯本區"] = "11"
	province_ch2code["CD:奎卢省"] = "KL"
	province_ch2code["TG:中部区"] = "C"
	province_ch2code["TG:卡拉区"] = "K"
	province_ch2code["MR:塔甘特省"] = "09"
	province_ch2code["MR:提里斯-宰穆爾省"] = "11"
	province_ch2code["MR:卜拉克納省"] = "05"
	province_ch2code["MR:因希里省"] = "12"
	province_ch2code["CF:洛巴耶省"] = "LB"
	province_ch2code["CF:瓦姆省"] = "AC"
	province_ch2code["BJ:阿黎博里省"] = "AL"
	province_ch2code["BJ:峽谷省"] = "DO"
	province_ch2code["GA:河口省"] = "1"
	province_ch2code["ST:普林西比岛"] = "P"
	province_ch2code["NE:迪法大区"] = "2"
	province_ch2code["ES:埃斯特雷马杜拉"] = "EX"
	province_ch2code["ES:卡斯蒂利亚-拉曼恰"] = "CM"
	province_ch2code["ES:加那利群岛"] = "CN"
	province_ch2code["IT:西西里岛"] = "82"
	province_ch2code["IT:卡拉布里亚"] = "78"
	province_ch2code["IT:巴斯利卡塔"] = "77"
	province_ch2code["IT:普利亚"] = "75"
	province_ch2code["IT:坎帕尼亚"] = "72"
	province_ch2code["DK:西兰大区"] = "85"
	province_ch2code["DK:南丹麦大区"] = "83"
	province_ch2code["DK:中日德兰大区"] = "82"
	province_ch2code["GB:英格兰"] = "ENG"
	province_ch2code["GB:苏格兰"] = "SCT"
	province_ch2code["CH:汝拉州"] = "JU"
	province_ch2code["SE:厄勒布鲁省"] = "T"
	province_ch2code["SE:西约塔兰省"] = "O"
	province_ch2code["SE:耶夫勒堡省"] = "X"
	province_ch2code["SJ:斯瓦尔巴群岛"] = "21"
	province_ch2code["NL:北布拉班特省"] = "NB"
	province_ch2code["AT:克恩顿州"] = "2"
	province_ch2code["AT:福拉尔贝格州"] = "8"
	province_ch2code["AT:布尔根兰州"] = "1"
	province_ch2code["DE:巴伐利亚"] = "BY"
	province_ch2code["DE:巴登-符腾堡"] = "BW"
	province_ch2code["DE:莱茵兰-普法尔茨"] = "RP"
	province_ch2code["DE:勃兰登堡"] = "BB"
	province_ch2code["DE:萨克森-安哈尔特"] = "ST"
	province_ch2code["DE:石勒苏益格-荷尔斯泰因"] = "SH"
	province_ch2code["DE:下萨克森"] = "NI"
	province_ch2code["DE:柏林"] = "BE"
	province_ch2code["IE:芒斯特省"] = "M"
	province_ch2code["FR:卢瓦尔河地区"] = "PDL"
	province_ch2code["FR:布列塔尼半岛"] = "BRE"
	province_ch2code["FR:普罗旺斯-阿尔卑斯-蓝色海岸"] = "PAC"
	province_ch2code["AD:卡尼略"] = "02"
	province_ch2code["AD:圣胡利娅-德洛里亚"] = "06"
	province_ch2code["AD:奥尔迪诺"] = "05"
	province_ch2code["LI:毛伦"] = "04"
	province_ch2code["LI:埃申"] = "02"
	province_ch2code["HU:佐洛州"] = "ZA"
	province_ch2code["HU:沃什州"] = "VA"
	province_ch2code["CZ:中波希米亚州"] = "20"
	province_ch2code["CZ:南波希米亚州"] = "31"
	province_ch2code["PL:西里西亚省"] = "24"
	province_ch2code["PL:大波兰省"] = "30"
	province_ch2code["PL:下西里西亚省"] = "02"
	province_ch2code["ES:卡斯蒂利亚-莱昂"] = "CL"
	province_ch2code["ES:加利西亚"] = "GA"
	province_ch2code["ES:坎塔布里亚"] = "CB"
	province_ch2code["ES:拉里奥哈"] = "RI"
	province_ch2code["IT:威尼托"] = "34"
	province_ch2code["IT:利古里亞"] = "42"
	province_ch2code["IT:拉齐奥"] = "62"
	province_ch2code["IT:皮埃蒙特"] = "21"
	province_ch2code["IT:托斯卡纳"] = "52"
	province_ch2code["IT:阿布鲁佐"] = "65"
	province_ch2code["IT:马尔凯"] = "57"
	province_ch2code["IT:莫利塞"] = "67"
	province_ch2code["IT:翁布里亚"] = "55"
	province_ch2code["SM:基埃萨努欧瓦"] = "02"
	province_ch2code["SI:伊德里亞"] = "036"
	province_ch2code["SI:斯洛文尼亞科尼采"] = "114"
	province_ch2code["SI:斯洛文尼亞比斯特里察"] = "113"
	province_ch2code["SI:伊利爾斯卡比斯特里察"] = "038"
	province_ch2code["SI:拉多夫利察"] = "102"
	province_ch2code["SI:科佩尔"] = "050"
	province_ch2code["SI:皮夫卡"] = "091"
	province_ch2code["SI:普图伊"] = "096"
	province_ch2code["SI:普雷瓦列"] = "175"
	province_ch2code["ME:莫伊科瓦茨"] = "11"
	province_ch2code["SI:马里博尔"] = "070"
	province_ch2code["SI:伊佐拉"] = "040"
	province_ch2code["ME:新海尔采格"] = "08"
	province_ch2code["SI:赫拉斯特尼克"] = "034"
	province_ch2code["SI:洛加泰茨"] = "064"
	province_ch2code["ME:比耶洛波列"] = "04"
	province_ch2code["SI:采爾克諾"] = "014"
	province_ch2code["AO:万博省"] = "HUA"
	province_ch2code["AO:本吉拉省"] = "BGU"
	province_ch2code["BR:伯南布哥州"] = "PE"
	province_ch2code["BR:托坎廷斯州"] = "TO"
	province_ch2code["BR:马拉尼昂州"] = "MA"
	province_ch2code["BR:帕拉州"] = "PA"
	province_ch2code["BR:亚马孙州"] = "AM"
	province_ch2code["GL:凯克卡塔"] = "QE"
	province_ch2code["GL:库雅雷克"] = "KU"
	province_ch2code["AR:布宜诺斯艾利斯"] = "C"
	province_ch2code["UY:塞罗拉尔戈省"] = "CL"
	province_ch2code["UY:卡内洛内斯省"] = "CA"
	province_ch2code["UY:阿蒂加斯省"] = "AR"
	province_ch2code["BR:米纳斯吉拉斯州"] = "MG"
	province_ch2code["BR:圣埃斯皮里图州"] = "ES"
	province_ch2code["BR:戈亚斯州"] = "GO"
	province_ch2code["BR:马托格罗索州"] = "MT"
	province_ch2code["BR:南马托格罗索州"] = "MS"
	province_ch2code["MX:塔毛利帕斯州"] = "TAM"
	province_ch2code["MX:墨西哥州"] = "MEX"
	province_ch2code["MX:瓦哈卡州"] = "OAX"
	province_ch2code["MX:普埃布拉州"] = "PUE"
	province_ch2code["MX:莫雷洛斯州"] = "MOR"
	province_ch2code["MX:尤卡坦州"] = "YUC"
	province_ch2code["MX:恰帕斯州"] = "CHP"
	province_ch2code["MX:墨西哥城市"] = "CMX"
	province_ch2code["MX:塔巴斯科州"] = "TAB"
	province_ch2code["MX:克雷塔羅州"] = "QUE"
	province_ch2code["MX:坎佩切州"] = "CAM"
	province_ch2code["VE:阿马库罗三角洲州"] = "Y"
	province_ch2code["PE:安卡什大区"] = "ANC"
	province_ch2code["PE:洛雷托大区"] = "LOR"
	province_ch2code["PE:拉利伯塔德大区"] = "LAL"
	province_ch2code["PE:乌卡亚利大区"] = "UCA"
	province_ch2code["PE:通贝斯大区"] = "TUM"
	province_ch2code["PE:兰巴耶克大区"] = "LAM"
	province_ch2code["PE:皮乌拉地区"] = "PIU"
	province_ch2code["PE:亚马孙大区"] = "AMA"
	province_ch2code["PA:恩贝拉-沃内安特区"] = "EM"
	province_ch2code["PA:库纳雅拉特区"] = "KY"
	province_ch2code["CL:圣地亚哥首都大区"] = "RM"
	province_ch2code["CL:塔拉帕卡大区"] = "TA"
	province_ch2code["PE:帕斯科大区"] = "PAS"
	province_ch2code["PE:阿雷基帕大区"] = "ARE"
	province_ch2code["PE:胡宁大区"] = "JUN"
	province_ch2code["PE:马德雷德迪奥斯大区"] = "MDD"
	province_ch2code["PE:塔克纳大区"] = "TAC"
	province_ch2code["PE:阿亚库乔大区"] = "AYA"
	province_ch2code["PE:伊卡大区"] = "ICA"
	province_ch2code["PE:普诺大区"] = "PUN"
	province_ch2code["MX:奇瓦瓦州"] = "CHH"
	province_ch2code["MX:哈利斯科州"] = "JAL"
	province_ch2code["MX:米却肯州"] = "MIC"
	province_ch2code["WS:艾加伊勒泰"] = "AL"
	province_ch2code["WS:阿图阿"] = "AT"
	province_ch2code["WS:图阿马萨加"] = "TU"
	province_ch2code["FJ:东部大区"] = "E"
	province_ch2code["US:德克萨斯州"] = "TX"
	province_ch2code["US:亚拉巴马州"] = "AL"
	province_ch2code["US:弗吉尼亚州"] = "VA"
	province_ch2code["US:阿肯色州"] = "AR"
	province_ch2code["US:特拉华州"] = "DE"
	province_ch2code["US:佛罗里达州"] = "FL"
	province_ch2code["US:乔治亚"] = "GA"
	province_ch2code["US:伊利诺伊州"] = "IL"
	province_ch2code["US:印第安纳州"] = "IN"
	province_ch2code["US:马里兰州"] = "MD"
	province_ch2code["US:肯塔基州"] = "KY"
	province_ch2code["US:密苏里州"] = "MO"
	province_ch2code["US:密西西比州"] = "MS"
	province_ch2code["US:北卡罗来纳州"] = "NC"
	province_ch2code["US:南卡罗来纳州"] = "SC"
	province_ch2code["US:田纳西州"] = "TN"
	province_ch2code["US:路易斯安那州"] = "LA"
	province_ch2code["US:新泽西州"] = "NJ"
	province_ch2code["US:俄亥俄州"] = "OH"
	province_ch2code["US:宾夕法尼亚州"] = "PA"
	province_ch2code["US:康乃狄克州"] = "CT"
	province_ch2code["US:艾奥瓦州"] = "IA"
	province_ch2code["US:缅因州"] = "ME"
	province_ch2code["US:密歇根州"] = "MI"
	province_ch2code["US:纽约州"] = "NY"
	province_ch2code["US:南达科他州"] = "SD"
	province_ch2code["US:威斯康辛州"] = "WI"
	province_ch2code["US:明尼苏达州"] = "MN"
	province_ch2code["US:北达科他州"] = "ND"
	province_ch2code["US:新罕布什尔州"] = "NH"
	province_ch2code["US:佛蒙特州"] = "VT"
	province_ch2code["US:加利福尼亚州"] = "CA"
	province_ch2code["US:新墨西哥州"] = "NM"
	province_ch2code["US:犹他州"] = "UT"
	province_ch2code["US:内华达州"] = "NV"
	province_ch2code["US:爱达荷州"] = "ID"
	province_ch2code["US:阿拉斯加州"] = "AK"
	province_ch2code["US:蒙大拿州"] = "MT"
	province_ch2code["US:俄勒冈州"] = "OR"
	province_ch2code["US:华盛顿州"] = "WA"
	province_ch2code["US:怀俄明州"] = "WY"
	province_ch2code["US:夏威夷州"] = "HI"
	province_ch2code["CA:不列颠哥伦比亚"] = "BC"
	province_ch2code["CA:艾伯塔"] = "AB"
	province_ch2code["CA:安大略"] = "ON"
	province_ch2code["CA:新斯科舍"] = "NS"
	province_ch2code["CA:曼尼托巴"] = "MB"
	province_ch2code["CA:新不倫瑞克"] = "NB"
	province_ch2code["CA:育空"] = "YT"
	province_ch2code["ES:休达"] = "CE"
	province_ch2code["MY:布城"] = "16"

	city_ch2en = make(map[string]string)
	city_en2code = make(map[string]string)
}

func continentCh2En(ch string) string {
	if val, ok := continents_ch2en[ch]; ok {
		return val
	}

	log.Printf("can not convert continent %s to english", ch)
	return ""
}

func continentCh2Code(ch string) string {
	if val, ok := continents_ch2code[ch]; ok {
		return val
	}

	log.Printf("can not convert continent %s to english", ch)
	return ""
}

func provinceCh2En(country_code string, ch string) string {
	key := country_code + ":" + ch
	if val, ok := province_ch2en[key]; ok {
		return val
	}

	log.Printf("can not convert province %s to english", ch)
	return ""
}

func provinceCh2Code(country_code string, ch string) string {
	key := country_code + ":" + ch
	if val, ok := province_ch2code[key]; ok {
		return val
	}

	log.Printf("can not convert province %s to code", ch)
	return ""
}

func cityCh2En(ch string) string {
	if val, ok := city_ch2en[ch]; ok {
		return val
	}

	log.Printf("can not convert city %s to english", ch)
	return ""
}

func cityEn2Code(en string) string {
	if val, ok := city_en2code[en]; ok {
		return val
	}

	log.Printf("can not convert city %s to english", en)
	return ""
}

func main() {

	qqip, err := qqzeng.New()
	if err != nil {
		return
	}

	qqip.GetAll()

	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "My-ASN-DB",
			RecordSize:   24,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range []string{"qqzeng.txt"} {
		fh, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		fileScanner := bufio.NewScanner(fh)
		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			line := fileScanner.Text()
			if len(line) == 0 {
				continue
			}

			row := strings.Split(line, "|")

			if len(row) != 15 {
				log.Fatalf("column number %d is not equal to 15, line: %v", len(row), line)
			}

			//61.134.201.0|61.134.211.255|1032243456|1032246271|亚洲|中国|山西|太原||联通|140100|China|CN|112.        549248|37.857014
			ip_start := row[0]
			//ip_end := row[1]
			ipnum_start, err := strconv.Atoi(row[2])
			if err != nil {
				log.Fatal("error: %s, line: %s", err, line)
			}

			ipnum_end, err := strconv.Atoi(row[3])
			if err != nil {
				log.Fatal("error: %s, line: %s", err, line)
			}

			if ipnum_end < ipnum_start {
				log.Fatal("ip_start should greater than ip_end, line: %s", err, line)
			}

			continentCh := row[4]
			country_ch := row[5]
			province_ch := row[6]
			city_ch := row[7]
			//district_ch := row[8]
			//isp := row[9]
			//zipcode := row[10]
			country_en := row[11]
			country_code := row[12]
			latitude, err := strconv.ParseFloat(row[13], 64)
			if err != nil {
				log.Fatal("error: %s, line: %s", err, line)
			}

			longitude, err := strconv.ParseFloat(row[14], 64)
			if err != nil {
				log.Fatal("error: %s, line: %s", err, line)
			}

			mask_len := int(math.Log2(float64(ipnum_end + 1 - ipnum_start)))

			_, network, err := net.ParseCIDR(ip_start + "/" + strconv.Itoa(mask_len))
			if err != nil {
				log.Fatal(err)
			}

			record := mmdbtype.Map{}

			if len(continentCh) > 0 {
				mmdbContinent := mmdbtype.Map{}
				names := mmdbtype.Map{"zh-CN": mmdbtype.String(continentCh)}
				continentEn := continentCh2En(continentCh)
				if len(continentEn) > 0 {
					continentCode := continentCh2Code(continentCh)
					mmdbContinent["code"] = mmdbtype.String(continentCode)
					names["en"] = mmdbtype.String(continentEn)
				}

				mmdbContinent["names"] = names
				record["continent"] = mmdbContinent
			}

			{
				mmdbCountry := mmdbtype.Map{}
				names := mmdbtype.Map{}
				if len(country_ch) > 0 {
					names["zh-CN"] = mmdbtype.String(country_ch)
				}

				if len(country_en) > 0 {
					names["en"] = mmdbtype.String(country_en)
				}

				mmdbCountry["names"] = names
				mmdbCountry["iso_code"] = mmdbtype.String(country_code)
				record["country"] = mmdbCountry
			}

			if len(province_ch) > 0 {
				mmdbProvice := mmdbtype.Map{}
				names := mmdbtype.Map{"zh-CN": mmdbtype.String(province_ch)}
				provinceEn := provinceCh2En(country_code, province_ch)
				if len(provinceEn) > 0 {
					names["en"] = mmdbtype.String(provinceEn)
				}
				provinceCode := provinceCh2Code(country_code, province_ch)
				if len(provinceCode) > 0 {
					mmdbProvice["iso_code"] = mmdbtype.String(provinceCode)
				}
				mmdbProvice["names"] = names
				record["subdivisions"] = mmdbProvice
			}

			if len(city_ch) > 0 {
				mmdbCity := mmdbtype.Map{}
				names := mmdbtype.Map{"zh-CN": mmdbtype.String(city_ch)}
				cityEn := cityCh2En(province_ch)
				if len(cityEn) > 0 {
					names["en"] = mmdbtype.String(cityEn)
				}

				mmdbCity["names"] = names
				record["city"] = mmdbCity
			}

			if latitude != 0 && longitude != 0 {
				record["location"] = mmdbtype.Map{
					"latitude":  mmdbtype.Float64(latitude),
					"longitude": mmdbtype.Float64(longitude),
				}
			}

			/*
				if asn != 0 {
					record["autonomous_system_number"] = mmdbtype.Uint32(asn)
				}

				if row[2] != "" {
					record["autonomous_system_organization"] = mmdbtype.String(row[2])
				}
			*/
			err = writer.Insert(network, record)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fh, err := os.Create("out.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
