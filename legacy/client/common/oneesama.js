class OneeSama {

	readableUTCTime(d, seconds) {
		let html = pad(d.getUTCDate()) + ' '
			+ this.lang.year[d.getUTCMonth()] + ' '
			+ d.getUTCFullYear()
			+ `(${this.lang.week[d.getUTCDay()]})`
			+`${pad(d.getUTCHours())}:${pad(d.getUTCMinutes())}`;
		if (seconds)
			html += `:${pad(d.getUTCSeconds())}`;
		html += ' UTC';
		return html;
	}


	expansionLinks(num) {
		return parseHTML
			`<span class="act expansionLinks">
				<a href="${num}" class="history">
					${this.lang.expand}
				</a>
				] [
				<a href="${num}?last=${this.lastN}" class="history">
					${this.lang.last} ${this.lastN}
				</a>
			</span>`;
	}

	asideLink(inner, href, cls, innerCls) {
		return parseHTML
			`<aside class="act glass ${cls}">
				<a ${href && `href="${href}"`}
					${innerCls && ` class="${innerCls}"`}
				>
					${this.lang[inner] || inner}
				</a>
			</aside>`
	}

	replyBox() {
		return this.asideLink('reply', null, 'posting');
	}
}
