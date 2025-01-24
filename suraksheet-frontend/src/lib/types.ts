export type Bin = {
	id: number
	name: string
	owner: number
	createdAt: Date
}

export type Document = {
	id: number
	name: string
	referenceName: string
	bin: number
	url: string|null
	extract: string
	createdAt: Date
}