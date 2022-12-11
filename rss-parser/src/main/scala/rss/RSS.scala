package rss

import zio._
import scala.io.Source
import scala.util.Using

case class RSS(title: String, url: String, rss: String)

object RSS {
  private def validateRSS(rss: scala.xml.Elem): IO[String, scala.xml.Elem] = rss match {
    case <rss>{elems @ _ *}</rss> =>
      val channel = for (elem @ <channel>{_ * }</channel> <- elems) yield elem
      if (channel.length != 1) ZIO.fail("no only one channel tag")
      else if ((channel.head \ "title").length + (channel.head \ "description").length == 2)
        ZIO.succeed(rss)
      else ZIO.fail("title and description tags occurred not only once")
    case _ => ZIO.fail("no rss tag")
  }

  def fromUrl(url: String) = ???

  def fromRSSUrl(url: String): IO[Serializable, RSS] = ZIO.fromTry(Using(Source.fromURL(url)) {
    source => scala.xml.XML.loadString(source.mkString)
  })
    //.flatMap(x => Console.printLine(x).map(_ => x))
    .foldZIO (
      _    => ZIO.fail("not rss link"),
      elem => validateRSS(elem).map(rss => RSS((rss \ "channel" \ "title").text, url, url))
    )

}