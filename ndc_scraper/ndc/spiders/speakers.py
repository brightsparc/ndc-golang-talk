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
            next_page = response.urljoin(speaker)
            yield scrapy.Request(next_page, callback=self.parse_speaker)

    def strip_not_empty(self, text):
        if text: return text.strip()

    def join_not_empty(self, text):
        if text: return re.sub('\s+', ' ', ' '.join(text))

    def parse_speaker(self, response):
        item = {
            "conference": response.css("div.logo p::text").re("\w+")[0],
            "name": self.strip_not_empty(response.css("section.masthead h1::text").extract_first()),
            "tagline": self.strip_not_empty(response.css("section.masthead h1 span::text").extract_first()),
            "image": response.urljoin(response.css("section.masthead img::attr(src)").extract_first()),
            "handle": self.strip_not_empty(response.css("section.masthead a::text").extract_first()),
            "preamble": self.join_not_empty(response.css("section.preamble p::text").extract()),
        }

        for talk in response.css("section.events li a::attr(href)").extract():
            next_page = response.urljoin(talk)
            yield scrapy.Request(next_page, callback=self.parse_talk, meta={"item": item})

    def parse_talk(self, response):
        item = response.meta["item"]
        print(response.url)
        tags = response.css("section.masthead div.tags a::text").extract()
        levels = [s[7:] for s in tags if s.startswith('Level: ')]
        if not levels:
            level = 'All levels'
        else:
            level = levels[0]
        content = {
            "url": response.url,
            "level": level,
            "title": self.strip_not_empty(response.css("section.masthead h1::text").extract_first()),
            "tags": [s for s in tags if not s.startswith('Level: ')], # don't include levels
            "preamble": self.join_not_empty(response.css("section.preamble p::text").extract()),
            "body": self.join_not_empty(response.css("section.body p::text").extract()),
        }
        if re.search("/workshop/", response.url):
            item["workshop"] = content
        else:
            item["talk"] = content
        return item
