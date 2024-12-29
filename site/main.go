package main

import (
	"syscall/js"
	"time"

	"github.com/hdck007/goat/goat"
)

// HeroSection component with glowing gradient
func HeroSection(props goat.Props) goat.Block {
	return goat.BlockElement(func(p goat.Props) goat.Vnode {
		return goat.CreateVirtualElements("div",
			goat.Props{"class": "bg-zinc-900 text-white py-32 relative overflow-hidden"},
			// Glow effects
			goat.CreateVirtualElements("div",
				goat.Props{"class": "absolute h-[60%] inset-0 bg-gradient-to-r from-blue-500/20 to-purple-500/20 blur-3xl"},
			),
			goat.CreateVirtualElements("div",
				goat.Props{"class": "container mx-auto px-6 text-center relative z-10"},
				goat.CreateVirtualElements("h1",
					goat.Props{"class": "text-6xl font-extrabold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-blue-300 to-purple-300"},
					goat.CreateVirtualElements("text", nil,
						&goat.TextOrElement{
							StringValue: "Empower Your Apps with Goat 🐐",
							Element:     nil,
						},
					),
				),
				goat.CreateVirtualElements("p",
					goat.Props{"class": "text-xl text-blue-100/90 mb-8 max-w-2xl mx-auto"},
					goat.CreateVirtualElements("text", nil,
						&goat.TextOrElement{
							StringValue: "Build blazing-fast, modern web apps with the power of Golang and WASM.",
							Element:     nil,
						},
					),
				),
			),
		)
	}, props)()
}

// Counter component with glowing border
func Counter(props goat.Props) goat.Block {
	return goat.BlockElement(func(p goat.Props) goat.Vnode {
		return goat.CreateVirtualElements("div",
			goat.Props{"class": "relative group"},
			// Glow effect
			goat.CreateVirtualElements("div",
				goat.Props{"class": "absolute -inset-0.5 bg-zinc-900 pt-8 rounded-lg blur  group-hover:opacity-75 transition duration-300"},
			),
			goat.CreateVirtualElements("div",
				goat.Props{"class": "relative bg-zinc-900 p-8 py-0 rounded-lg"},
				goat.CreateVirtualElements("div",
					goat.Props{"class": "text-center"},
					goat.CreateVirtualElements("h2",
						goat.Props{"class": "text-3xl font-bold text-white mb-4"},
						goat.CreateVirtualElements("text", nil,
							&goat.TextOrElement{
								StringValue: "Counter: ",
								Element:     nil,
							},
						),
						goat.Get("count"),
					),
					goat.CreateVirtualElements("p",
						goat.Props{"class": "text-blue-200"},
						goat.CreateVirtualElements("text", nil,
							&goat.TextOrElement{
								StringValue: "Real-time updates with VDOM",
								Element:     nil,
							},
						),
					),
				),
			),
		)
	}, props)()
}

// Features section with glowing cards
func FeaturesSection(props goat.Props) goat.Block {
	return goat.BlockElement(func(p goat.Props) goat.Vnode {
		return goat.CreateVirtualElements("div",
			goat.Props{"class": "bg-zinc-900 py-16"}, // Section background
			goat.CreateVirtualElements("div",
				goat.Props{"class": "container mx-auto px-6"},
				goat.CreateVirtualElements("div",
					goat.Props{"class": "max-w-screen-md mb-8 lg:mb-16"},

					// Removed the description here
				),
				goat.CreateVirtualElements("div",
					goat.Props{"class": "space-y-8 md:grid md:grid-cols-2 lg:grid-cols-3 md:gap-12 md:space-y-0"},

					// Feature 1: WebAssembly Power 🕹️
					goat.CreateVirtualElements("div",
						goat.Props{"class": "text-center relative group"},
						goat.CreateVirtualElements("div",
							goat.Props{"class": "absolute -inset-0.5 bg-gradient-to-r from-blue-500 to-purple-500 rounded-lg blur opacity-0 group-hover:opacity-50 transition duration-300"}, // Gradient outline hidden by default, shown on hover
						),
						goat.CreateVirtualElements("div",
							goat.Props{"class": "relative p-6 bg-zinc-900 dark:bg-gray-800 rounded-lg shadow-lg"}, // Card background updated to bg-zinc-900
							goat.CreateVirtualElements("h3",
								goat.Props{"class": "text-xl font-semibold text-white mb-2"}, // Text color adjusted
								goat.CreateVirtualElements("text", nil,
									&goat.TextOrElement{
										StringValue: "WebAssembly Power 🕹️",
										Element:     nil,
									},
								),
							),
							goat.CreateVirtualElements("p",
								goat.Props{"class": "text-blue-200 min-h-[3rem] line-clamp-2"}, // Description forces at least 2 lines
								goat.CreateVirtualElements("text", nil,
									&goat.TextOrElement{
										StringValue: "Compiled to WASM for native-like performance",
										Element:     nil,
									},
								),
							),
						),
					),

					// Feature 2: Fast Updates ⚡
					goat.CreateVirtualElements("div",
						goat.Props{"class": "text-center relative group"},
						goat.CreateVirtualElements("div",
							goat.Props{"class": "absolute -inset-0.5 bg-gradient-to-r from-blue-500 to-purple-500 rounded-lg blur opacity-0 group-hover:opacity-50 transition duration-300"}, // Gradient outline hidden by default, shown on hover
						),
						goat.CreateVirtualElements("div",
							goat.Props{"class": "relative p-6 bg-zinc-900 dark:bg-gray-800 rounded-lg shadow-lg"}, // Card background updated to bg-zinc-900
							goat.CreateVirtualElements("h3",
								goat.Props{"class": "text-xl font-semibold text-white mb-2"}, // Text color adjusted
								goat.CreateVirtualElements("text", nil,
									&goat.TextOrElement{
										StringValue: "Fast Updates ⚡",
										Element:     nil,
									},
								),
							),
							goat.CreateVirtualElements("p",
								goat.Props{"class": "text-blue-200 min-h-[3rem] line-clamp-2"}, // Description forces at least 2 lines
								goat.CreateVirtualElements("text", nil,
									&goat.TextOrElement{
										StringValue: "Lightning-fast DOM updates with minimal operations",
										Element:     nil,
									},
								),
							),
						),
					),

					// Feature 3: Efficient Updates 🔧
					goat.CreateVirtualElements("div",
						goat.Props{"class": "text-center relative group"},
						goat.CreateVirtualElements("div",
							goat.Props{"class": "absolute -inset-0.5 bg-gradient-to-r from-blue-500 to-purple-500 rounded-lg blur opacity-0 group-hover:opacity-50 transition duration-300"}, // Gradient outline hidden by default, shown on hover
						),
						goat.CreateVirtualElements("div",
							goat.Props{"class": "relative p-6 bg-zinc-900 dark:bg-gray-800 rounded-lg shadow-lg"}, // Card background updated to bg-zinc-900
							goat.CreateVirtualElements("h3",
								goat.Props{"class": "text-xl font-semibold text-white mb-2"}, // Text color adjusted
								goat.CreateVirtualElements("text", nil,
									&goat.TextOrElement{
										StringValue: "Efficient Updates 🔧",
										Element:     nil,
									},
								),
							),
							goat.CreateVirtualElements("p",
								goat.Props{"class": "text-blue-200 min-h-[3rem] line-clamp-2"}, // Description forces at least 2 lines
								goat.CreateVirtualElements("text", nil,
									&goat.TextOrElement{
										StringValue: "Smart diffing algorithm for optimal performance",
										Element:     nil,
									},
								),
							),
						),
					),
				),
			),
		)
	}, props)()
}

// FooterSection component with GitHub and Live Example links
func FooterSection(props goat.Props) goat.Block {
	return goat.BlockElement(func(p goat.Props) goat.Vnode {
		return goat.CreateVirtualElements("div",
			goat.Props{"class": "absolute bottom-4 right-4 text-white flex items-center space-x-4 z-20"},
			// GitHub Link with Icon
			goat.CreateVirtualElements("a",
				goat.Props{
					"href":   "https://github.com/hdck007/goat",
					"target": "_blank",
					"rel":    "noopener noreferrer",
					"class":  "flex items-center space-x-2 hover:text-blue-300",
				},
				goat.CreateVirtualElements("img",
					goat.Props{
						"src":   "/github-icon.png",
						"alt":   "GitHub Icon",
						"class": "w-6 h-6",
					},
				),
				goat.CreateVirtualElements("span",
					goat.Props{"class": "text-sm"},
					goat.CreateVirtualElements("text", nil,
						&goat.TextOrElement{
							StringValue: "GitHub",
						},
					),
				),
			),

			// Live Example Link
			goat.CreateVirtualElements("a",
				goat.Props{
					"href":   "/todoapp",
					"target": "_blank",
					"rel":    "noopener noreferrer",
					"class":  "text-sm hover:text-blue-300",
				},
				goat.CreateVirtualElements("text", nil,
					&goat.TextOrElement{
						StringValue: "Live Example",
					},
				),
			),
		)
	}, props)()
}

func App(props goat.Props) goat.Block {
	heroSection := HeroSection(goat.Props{})
	counterComponent := Counter(goat.Props{
		"count": props["count"],
	})
	featuresSection := FeaturesSection(goat.Props{})
	footer := FooterSection(goat.Props{})

	className := "min-h-screen bg-zinc-900 overflow-hidden"
	derivedProps := goat.Props{
		"heroSection":      heroSection,
		"counterComponent": counterComponent,
		"featuresSection":  featuresSection,
		"footer":           footer,
		"className":        className,
	}

	return goat.BlockElement(func(p goat.Props) goat.Vnode {
		return goat.CreateVirtualElements("div",
			goat.Props{"class": goat.Get("className")},
			goat.Get("heroSection"),
			goat.Get("counterComponent"),
			goat.Get("featuresSection"),
			goat.Get("footer"),
		)
	}, derivedProps)()
}

func main() {

	done := make(chan struct{}, 0)

	count := 0
	component := App(goat.Props{
		"count": count,
	})
	root := js.Global().Get("document").Call("getElementById", "root")
	component.Mount(root)

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				count++
				newComponent := App(goat.Props{
					"count": count,
				})
				component.Patch(newComponent)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-done
}
