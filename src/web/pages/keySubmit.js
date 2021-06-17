import Head from 'next/head'

import Protector from '../components/ProtectorComponent'
import KeySubmit from '../components/offlineKeySubmit'

export default function Home() {
  return (
    <div>
      <Head>
        <title>RW CnC</title>
        <meta name="description" content="cool page brah" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Protector shouldBeLoggedIn={true} goTo={"/"}>
        <KeySubmit />
      </Protector>
    </div>
  )
}
