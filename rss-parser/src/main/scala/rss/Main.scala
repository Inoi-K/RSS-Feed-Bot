package rss

import zio._
import zhttp.http._
import zhttp.service.Server
import zio.json._

case class Link(url: String)
case class Validity(valid: Boolean)

object Link {
  implicit val encoder: JsonEncoder[Link] =
    DeriveJsonEncoder.gen[Link]

  implicit val decoder: JsonDecoder[Link] =
    DeriveJsonDecoder.gen[Link]
}

object Validity {
  implicit val encoder: JsonEncoder[Validity] =
    DeriveJsonEncoder.gen[Validity]

  implicit val decoder: JsonDecoder[Validity] =
    DeriveJsonDecoder.gen[Validity]
}

object LinkApp {
  def apply(): Http[Any, Throwable, Request, Response] = Http.collectZIO[Request] {
    case req @ Method.POST -> !! / "link" =>
      for {
        u <- req.body.asString.map(_.fromJson[Link])
        r <- u match {
          case Left(e) =>
            ZIO.debug(s"Failed to parse the input: $e")
              .as(Response.text(e).setStatus(Status.BadRequest))
          case Right(link) =>
            RSS.fromRSSUrl(link.url)
              .fold(_ => false, _ => true)
              .map(x => Response.text(Validity(x).toJson))
        }
      } yield r
  }
}

object Main extends ZIOAppDefault {

  val port = 8000

  def run = for {
    _ <- Console.printLine(s"starting server on http://localhost:$port")
    _ <- Server.start(port, LinkApp())
    //RSS.scala.fromRSSUrl("https://code.visualstudio.com/feed.xml") -- not valid example
    //rss <- RSS.fromRSSUrl("https://rss.art19.com/apology-line")
    //_   <- Console.printLine(rss)
  } yield ()

}
