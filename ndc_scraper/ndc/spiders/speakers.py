import re
import scrapy

class SpeakersSpider(scrapy.Spider):
    name = "speakers"
    start_urls = [
        'http://ndcsydney.com/speakers/',
        'http://ndcoslo.com/speakers/'
    ]

    def parse(self, response):
        for speaker in response.css("a.boxed-speaker::attr(href)").extract():
            next_page = response.urljoin(speaker).strip('/') # Strip trailing slash
            yield scrapy.Request(next_page, callback=self.parse_speaker)

    def strip_not_empty(self, text):
        if text: return text.strip()

    def join_not_empty(self, text):
        if text: return re.sub('\s+', ' ', ' '.join(text))

    def parse_speaker(self, response):
        item = {
            "conference": response.css("div.logo p::text").re("\w+")[0],
            "speaker": {
                "url": response.url,
                "name": self.strip_not_empty(response.css("section.masthead h1::text").extract_first()),
                "tagline": self.strip_not_empty(response.css("section.masthead h1 span::text").extract_first()),
                "imageUrl": response.urljoin(response.css("section.masthead img::attr(src)").extract_first()),
                "handle": self.strip_not_empty(response.css("section.masthead a::text").extract_first()),
                "preamble": self.join_not_empty(response.css("section.preamble p::text").extract()),
            }
        }

        for talk in response.css("section.events li a"):
            link = response.urljoin(talk.css("::attr(href)").extract_first())
            item["talk"] = {
                "url": link,
                "workshop": re.search("/workshop/", link) is not None,
                "title": talk.css("h2::text").extract_first(),
            }
            # Reqeust the talk, but don't filter as we want to get multiple talks from the different speakers
            yield scrapy.Request(link, callback=self.parse_talk, dont_filter=True, meta={"item": item})

    def parse_talk(self, response):
        item = response.meta["item"]
        talk = item["talk"]
        print(response.url)
        # Get speaker links, if more than one pass through the 1-based order
        speaker_links = [response.urljoin(link) for link in response.css("section.speakers li a::attr(href)").extract()]
        if speaker_links and len(speaker_links) > 0:
            try:
                talk["order"] = speaker_links.index(item["speaker"]["url"]) + 1
                print('Got speaker order', talk["order"], speaker_links)
            except ValueError as e:
                print("Error getting speaker order", e)
        # Get tags excluding level
        tags = response.css("section.masthead div.tags a::text").extract()
        talk["tags"] = [s for s in tags if not s.startswith('Level: ')]
        talk["preamble"] = self.join_not_empty(response.css("section.preamble p::text").extract())
        talk["body"] = self.join_not_empty(response.css("section.body p::text").extract())
        # Get day/time/venue/level attributes
        for detail in response.css("section.details p"):
            label = detail.css("p::text").extract_first().strip().lower()
            talk[label] = detail.css("span::text").extract_first().strip()
        return item
