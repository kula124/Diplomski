import Head from 'next/head'

import Protector from '../components/ProtectorComponent'
import V2 from '../components/v2Page'

export default function Home() {
  return (
    <div>
      <Head>
        <title>RW CnC</title>
        <meta name="description" content="v2 page" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Protector shouldBeLoggedIn={true} goTo={"/"}>
        <V2 />
      </Protector>
    </div>
  )
}