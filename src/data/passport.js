const mapsUrl = (name, address) =>
  `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(`${name} ${address}`)}`

const store = (name, address) => ({ name, address, href: mapsUrl(name, address) })

const yeliu = [
  store('魚村活海鮮', '新北市萬里區野柳里港東路74之6號'),
  store('野柳美觀園飯店', '新北市萬里區野柳里港東路156號'),
  store('建香四姊妹海產店', '新北市萬里區野柳里港東路156號'),
  store('柏榕活海鮮', '新北市萬里區野柳里港東151-2號'),
  store('錦園漁夫料理', '新北市萬里區野柳里港東路159號'),
  store('珠海活海鮮餐廳', '新北市萬里區野柳里港東路162-13號'),
  store('女皇餐廳', '新北市萬里區野柳里港東路163號'),
  store('三葉活海鮮', '新北市萬里區野柳里港東路74號之16'),
  store('望海亭三嫂的店', '新北市野柳里港東路162-15號'),
]

const wanli = [
  store('金湧泉SPA溫泉會館', '新北市萬里區萬里加投213之3號'),
  store('沐舍溫泉渡假酒店', '新北市萬里區加投166-1號'),
  store('呷尚鱻餐廳／尚鱻漁舖', '基隆市基金三路67-1號'),
  store('魚多多海鮮料理', '新北市萬里區美崙路32號'),
  store('湄公河南洋美食', '新北市萬里區瑪鋉路244號'),
]

const eastCoast = [
  store('魚美村活海鮮', '新北市萬里區野柳里東澳路65號'),
  store('阿達活海產餐廳', '新北市石門區富基里楓林17-10號'),
  store('姊妹海景餐廳', '新北市石門區富基村楓林15-16號'),
  store('春金活海產', '新北市石門區富基里楓林路17-8號'),
  store('黑白毛海鮮', '新北市貢寮區仁和路51號'),
  store('嘉邑海鮮小館', '新北市貢寮區仁和路2號'),
  store('富士海鮮餐廳', '新北市貢寮區福隆里東興街8號'),
  store('魚佳餚餐館（東興街）', '新北市貢寮區福隆里東興街36號'),
  store('魚佳餚餐館（仁和路）', '新北市貢寮區仁里里仁和路186號'),
]

const guiHou = [
  store('小漁村海鮮店', '新北市萬里區龜吼里漁澳路63號'),
  store('珍鼎豐海鮮景觀餐廳', '新北市萬里區龜吼里漁澳路2號2樓'),
  store('漁家鄉海鮮餐廳', '新北市萬里區龜吼里漁澳路83號1樓'),
  store('金翡翠餐廳', '新北市萬里區龜吼里漁澳路16之3號'),
  store('富港海鮮', '新北市萬里區龜吼里漁澳路64-9號'),
  store('阿嬌海鮮館', '新北市萬里區龜吼里漁澳路80號'),
  store('北海漁港', '新北市萬里區龜吼里漁澳路16號'),
  store('巧晏漁坊', '新北市萬里區龜吼里漁澳路62號之1'),
  store('漁莊本港海鮮', '新北市萬里區龜吼里漁澳路64號之2'),
  store('麗鮮料理', '龜吼漁夫市集二樓攤號7'),
  store('珍海活海鮮', '新北市萬里區龜吼里漁澳路62-2號'),
  store('小微海產餐廳', '新北市萬里區龜吼里漁澳路16之1號'),
  store('一品鮮萬里蟹漁夫料理', '新北市萬里區龜吼里漁澳路65號'),
  store('津鮮閣海鮮餐廳', '新北市萬里區龜吼里漁澳路16之2號'),
  store('99蟹老闆', '新北市萬里區龜吼里石角路53號'),
  store('明發料理坊', '新北市萬里區漁澳路79-2號'),
  store('蟹港海鮮', '新北市萬里區漁澳路64之5號'),
]

const shenAo = [
  store('深澳港平價海鮮樓', '新北市瑞芳區深澳路16-7號'),
  store('海園活海鮮', '新北市瑞芳區鼻頭里鼻頭路245號'),
  store('88號水碼頭', '新北市金山區豐漁里民生路88號'),
  store('醉船長日式海鮮餐廳', '新北市金山區環金路200號'),
  store('田中芳園', '新北市金山區清水路53號'),
  store('花義思廚坊', '新北市金山區環金路286號'),
]

const shimen = [store('英芳飯店', '新北市石門區石門里中央路9號之5')]
const sanzhi = [
  store('八八八海鮮美食館', '新北市三芝區中興街2段30號之一'),
  store('鄉園本港海鮮餐廳', '新北市三芝區淡金路一段4號之1'),
]
const danshui = [
  store('魚藏海鮮宴會廣場', '新北市淡水區觀海路201號2樓'),
  store('大胖活海產', '新北市淡水區漁人碼頭商店街'),
]
const stirFry = [
  store('三和大盤熱炒三重店', '新北市三重區自強路三段80號'),
  store('不仔的店', '新北市新店區三民路135號'),
  store('開喜閣再來海鮮餐廳', '新北市新店區北新路一段88巷1號'),
  store('喝桶海啤酒屋海鮮餐', '新北市板橋區民生路三段309號'),
  store('海釣族真味園', '新北市板橋區文化路二段126號'),
  store('珠雞城', '新北市樹林區保安街一段85號'),
  store('新方海鮮宴會館', '新北市五股區成泰路一段194-2號'),
]

const rows = (stores, { left, top, width, step, height = 3.2 }) =>
  stores.map((item, index) => ({ ...item, left, top: top + (step * index), width, height }))

export const passportDesktopSlides = [
  {
    src: '/assets/pages/passport-desktop-1.jpg',
    alt: '海派護照合作亮點餐廳第 1 頁',
    links: [
      ...rows(yeliu, { left: 7.1, top: 29.6, width: 13.3, step: 3.72 }),
      ...rows(wanli, { left: 53.4, top: 29.6, width: 14.7, step: 3.72 }),
      ...rows(eastCoast, { left: 53.4, top: 60.0, width: 14.7, step: 3.72 }),
    ],
  },
  {
    src: '/assets/pages/passport-desktop-2.jpg',
    alt: '海派護照合作亮點餐廳第 2 頁',
    links: [
      ...rows(guiHou, { left: 6.1, top: 29.2, width: 14.7, step: 3.72 }),
      ...rows(shenAo, { left: 54.1, top: 32.2, width: 16.8, step: 4.0 }),
      ...rows(shimen, { left: 54.1, top: 68.3, width: 13.5, step: 4.0 }),
    ],
  },
  {
    src: '/assets/pages/passport-desktop-3.jpg',
    alt: '海派護照合作亮點餐廳第 3 頁',
    links: [
      ...rows(sanzhi, { left: 7.6, top: 28.3, width: 15.1, step: 4.0 }),
      ...rows(danshui, { left: 53.5, top: 28.3, width: 16.2, step: 4.0 }),
      ...rows(stirFry, { left: 29.8, top: 61.8, width: 15.1, step: 4.6 }),
    ],
  },
]

// The supplied mobile pages 2 and 3 contain the same restaurants. Page 3 is the
// clearer, table-based revision, so it is used as the second page without duplication.
export const passportMobileSlides = [
  {
    src: '/assets/pages/mobile/passport-mobile-1.jpg',
    alt: '海派護照合作亮點餐廳手機版第 1 頁',
    links: [
      ...rows(yeliu, { left: 13.3, top: 11.0, width: 20.5, step: 1.74, height: 1.7 }),
      ...rows(wanli, { left: 13.3, top: 31.0, width: 23.5, step: 1.82, height: 1.7 }),
      ...rows(eastCoast, { left: 13.3, top: 44.5, width: 21.2, step: 1.81, height: 1.7 }),
      ...rows(guiHou, { left: 13.3, top: 65.0, width: 24.5, step: 1.75, height: 1.65 }),
    ],
  },
  {
    src: '/assets/pages/mobile/passport-mobile-2.jpg',
    alt: '海派護照合作亮點餐廳手機版第 2 頁',
    links: [
      ...rows(shenAo, { left: 9.2, top: 16.1, width: 25.0, step: 2.45, height: 2.25 }),
      ...rows(shimen, { left: 9.2, top: 39.1, width: 19.0, step: 2.45, height: 2.25 }),
      ...rows(sanzhi, { left: 9.2, top: 48.6, width: 24.5, step: 2.55, height: 2.3 }),
      ...rows(danshui, { left: 9.2, top: 62.0, width: 25.0, step: 2.55, height: 2.3 }),
      ...rows(stirFry, { left: 9.2, top: 78.6, width: 26.3, step: 2.48, height: 2.3 }),
    ],
  },
]
