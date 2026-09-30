const mapsUrl = (name, address) =>
  `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(`${name} ${address}`)}`

const store = (name, address) => ({ name, address, href: mapsUrl(name, address) })

const yeliu = [
  store('魚村活海鮮', '新北市萬里區野柳里港東路74之6號'),
  store('野柳美觀園飯店', '新北市萬里區野柳里港東路156號'),
  store('建香四姊妹海產店', '新北市萬里區野柳里港東路160號'),
  store('錦園漁夫料理', '新北市萬里區野柳里港東路159號'),
  store('珠海活海鮮餐廳', '新北市萬里區野柳里港東路162-13號'),
  store('三葉活海鮮', '新北市萬里區野柳里港東路74號之16'),
]

const wanli = [
  store('金湧泉SPA溫泉會館', '新北市萬里區萬里加投213之3號'),
  store('沐舍溫泉渡假酒店', '新北市萬里區加投166-1號'),
]

const fuji = [
  store('阿達活海產餐廳', '新北市石門區富基里楓林17-10號'),
  store('姊妹海景餐廳', '新北市石門區富基村楓林15-16號'),
  store('春金活海產', '新北市石門區富基里楓林路17-8號'),
]

const guihou = [
  store('小漁村海鮮店', '新北市萬里區龜吼里漁澳路63號'),
  store('金翡翠餐廳', '新北市萬里區龜吼里漁澳路16之3號'),
  store('富港海鮮', '新北市萬里區龜吼里漁澳路64-9號'),
  store('阿嬌海鮮館', '新北市萬里區龜吼里漁澳路80號'),
  store('北海漁港', '新北市萬里區龜吼里漁澳路16號'),
  store('巧晏漁坊', '新北市萬里區龜吼里漁澳路62號之1'),
  store('漁莊本港海鮮', '新北市萬里區龜吼里漁澳路64號之2'),
  store('珍海活海鮮', '新北市萬里區龜吼里漁澳路62-2號'),
  store('一品鮮萬里蟹漁夫料理', '新北市萬里區龜吼里漁澳路65號'),
  store('津鮮閣海鮮餐廳', '新北市萬里區龜吼里漁澳路16之2號'),
  store('99蟹老闆', '新北市萬里區龜吼里石角路53號'),
  store('明發料理坊', '新北市萬里區漁澳路79-2號'),
  store('蟹港海鮮', '新北市萬里區漁澳路64之5號'),
]

const aodiFulong = [
  store('黑白毛海鮮', '新北市貢寮區仁和路51號'),
  store('嘉邑海鮮小館', '新北市貢寮區仁和路2號'),
  store('魚佳餚餐館', '新北市貢寮區福隆里東興街36號'),
]

const shenao = [store('深澳港平價海鮮樓', '新北市瑞芳區深澳路16-7號')]
const jinshan = [store('88號水碼頭', '新北市金山區豐漁里民生路88號')]
const shimen = [store('英芳飯店', '新北市石門區石門里中央路9號之5')]

const danshui = [
  store('魚藏海鮮宴會廣場', '新北市淡水區觀海路201號2樓'),
  store('大胖活海產', '新北市淡水區漁人碼頭商店街'),
]

const cityRestaurants = [
  store('珠雞城', '新北市樹林區保安街一段85號'),
  store('新方海鮮宴會館', '新北市五股區成泰路一段194-2號'),
]

const featuredRestaurants = [
  store('台北板橋馥華艾美酒店', '新北市板橋區中山路一段220號'),
  store('福容大飯店福隆', '新北市貢寮區福隆街41號'),
  store('新北北海溫泉洲際酒店', '新北市金山區環金路210號'),
  store('望月樓', '新北市板橋區新站路16號48樓'),
  store('香宮', '臺北市大安區敦化南路二段201號6樓'),
  store('上海醉月樓', '臺北市大安區敦化南路二段201號39樓'),
  store('JUNTO 同', '台北市中正區館前路22號正旅館藍2樓'),
]

const rows = (stores, { left, top, width, step, height = 3.1 }) =>
  stores.map((item, index) => ({ ...item, left, top: top + (step * index), width, height }))

export const passportDesktopSlides = [
  {
    src: '/assets/pages/passport-20260930-1.jpg',
    alt: '海派護照合作亮點餐廳第 1 頁',
    links: [
      ...rows(yeliu, { left: 7.2, top: 29.0, width: 15.2, step: 3.62 }),
      ...rows(wanli, { left: 7.2, top: 63.0, width: 16.5, step: 3.62 }),
      ...rows(fuji, { left: 7.2, top: 84.0, width: 15.6, step: 3.62 }),
      ...rows(guihou, { left: 53.5, top: 29.0, width: 16.0, step: 3.62 }),
    ],
  },
  {
    src: '/assets/pages/passport-20260930-2.jpg',
    alt: '海派護照合作亮點餐廳第 2 頁',
    links: [
      ...rows(aodiFulong, { left: 8.3, top: 29.0, width: 14.8, step: 3.72 }),
      ...rows(shenao, { left: 8.3, top: 50.0, width: 16.0, step: 3.72 }),
      ...rows(jinshan, { left: 8.3, top: 63.8, width: 14.0, step: 3.72 }),
      ...rows(shimen, { left: 8.3, top: 77.2, width: 13.0, step: 3.72 }),
      ...rows(danshui, { left: 8.3, top: 90.6, width: 16.0, step: 3.72 }),
      ...rows(cityRestaurants, { left: 55.0, top: 29.0, width: 16.0, step: 3.72 }),
      ...rows(featuredRestaurants, { left: 55.0, top: 49.1, width: 18.5, step: 3.72 }),
    ],
  },
]

export const passportMobileSlides = [
  {
    src: '/assets/pages/mobile/passport-20260930-1.jpg',
    alt: '海派護照合作亮點餐廳手機版第 1 頁',
    links: [
      ...rows(yeliu, { left: 12.4, top: 12.3, width: 24.0, step: 2.05, height: 1.85 }),
      ...rows(wanli, { left: 12.4, top: 31.0, width: 27.0, step: 2.05, height: 1.85 }),
      ...rows(guihou, { left: 12.4, top: 41.3, width: 27.0, step: 2.04, height: 1.85 }),
      ...rows(fuji, { left: 12.4, top: 74.4, width: 25.0, step: 2.05, height: 1.85 }),
      ...rows(aodiFulong, { left: 12.4, top: 85.1, width: 24.0, step: 2.05, height: 1.85 }),
    ],
  },
  {
    src: '/assets/pages/mobile/passport-20260930-2.jpg',
    alt: '海派護照合作亮點餐廳手機版第 2 頁',
    links: [
      ...rows(shenao, { left: 12.4, top: 12.2, width: 27.0, step: 2.1, height: 1.9 }),
      ...rows(jinshan, { left: 12.4, top: 21.1, width: 23.0, step: 2.1, height: 1.9 }),
      ...rows(shimen, { left: 12.4, top: 29.3, width: 21.0, step: 2.1, height: 1.9 }),
      ...rows(danshui, { left: 12.4, top: 38.0, width: 28.0, step: 2.15, height: 1.9 }),
      ...rows(cityRestaurants, { left: 12.4, top: 52.4, width: 28.0, step: 2.15, height: 1.9 }),
      ...rows(featuredRestaurants, { left: 12.4, top: 64.2, width: 31.0, step: 2.18, height: 1.95 }),
    ],
  },
]
