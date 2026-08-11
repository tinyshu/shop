import config from '@/config/config.js'

export const DEFAULT_IMAGE = '/static/nopicture.jpg'

export function resolveImageUrl(url) {
	if (!url) return DEFAULT_IMAGE
	if (url.slice(0, 4) === 'http') return url
	return config.baseUrl + '/' + url.replace(/^\//, '')
}

export function onImageError(item, key = 'imgUrl') {
	if (item && item[key] !== DEFAULT_IMAGE) {
		item[key] = DEFAULT_IMAGE
	}
}
