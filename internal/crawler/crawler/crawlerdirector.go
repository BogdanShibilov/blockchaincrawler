package crawler

import "context"

type director struct {
	crawlerBuilder CrawlerBuilder
}

func NewDirector(b CrawlerBuilder) *director {
	return &director{
		crawlerBuilder: b,
	}
}

func (d *director) BuildCrawler(ctx context.Context) (Crawler, error) {
	err := d.crawlerBuilder.setClient()
	if err != nil {
		return nil, err
	}
	d.crawlerBuilder.setBlocksCh()
	d.crawlerBuilder.setHeadersCh()
	d.crawlerBuilder.setErrorsCh()
	err = d.crawlerBuilder.setSubcription(ctx)
	if err != nil {
		return nil, err
	}
	return d.crawlerBuilder.buildCrawler(), nil
}
