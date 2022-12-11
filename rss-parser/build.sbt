ThisBuild / version := "0.1.0-SNAPSHOT"

ThisBuild / scalaVersion := "2.13.8"

lazy val root = (project in file("."))
  .settings(
    name := "rss-parser"
  )

libraryDependencies ++= Seq(
  "dev.zio" %% "zio" % "2.0.2",
  "dev.zio" %% "zio-json" % "0.3.0-RC8",
  "io.d11" %% "zhttp" % "2.0.0-RC11",
  "org.scala-lang.modules" %% "scala-xml" % "2.1.0"
)
