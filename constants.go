package syringe

type Lifecycle uint8

const Transient Lifecycle = 1
const Singleton Lifecycle = 2

const injectorTag = "di"
