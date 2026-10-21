package data

// English/international reference data for dummy documents. Every slice is
// package-level and must not be modified at runtime.

var englishData = dataset{
	firstNames: []string{
		"John", "Emily", "Michael", "Sarah", "David", "Jessica",
		"Daniel", "Laura", "James", "Olivia", "Robert", "Emma",
	},
	lastNames: []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Miller",
		"Davis", "Wilson", "Anderson", "Taylor", "Thomas", "Moore",
	},
	cities: []string{
		"New York", "London", "Toronto", "Sydney", "Chicago", "Boston",
		"Seattle", "Austin", "Melbourne", "Dublin", "San Francisco", "Vancouver",
	},
	streets: []string{
		"Main Street", "Oak Avenue", "Maple Drive", "Cedar Lane",
		"Park Road", "Lake View", "Hill Street", "Washington Avenue",
		"Pine Street", "River Road", "Sunset Boulevard", "Chestnut Street",
	},
	companies: []string{
		"Northwind Inc.", "Acme LLC", "Globex Corporation", "Initech LLC",
		"Umbrella Group", "Stark Industries", "Wayne Enterprises", "Soylent Corp.",
		"Hooli Inc.", "Vandelay Industries", "Pied Piper LLC", "Massive Dynamic",
	},
	goods: []string{
		"A4 Paper", "Printer Ink", "Gel Pen", "Plastic Folder",
		"Planner", "Calculator", "Stapler", "Whiteboard Marker",
		"32GB Flash Drive", "HDMI Cable", "Office Chair", "Folding Table",
	},
	units: []string{
		"pcs", "box", "ream", "unit", "set", "dozen",
		"bag", "meter", "liter", "kg", "pack", "sheet",
	},
	positions: []string{
		"Director", "Finance Manager", "General Affairs Head", "Administrative Staff",
		"Secretary", "Treasurer", "Operations Manager", "Warehouse Head",
		"Supervisor", "Coordinator", "Administrator", "Owner",
	},
	purposes: []string{
		"Consulting service payment", "Office supplies purchase",
		"Room rental payment", "Network maintenance service",
		"Computer equipment purchase", "Employee training cost",
		"Design service payment", "Stationery procurement",
		"Business travel cost", "Electricity bill payment",
		"Document printing service", "Raw material purchase",
	},
	deliveryNotes: []string{
		"Good", "As ordered", "No defects", "Packed neatly",
		"Fragile", "Ship today", "Checked", "Seal intact",
		"Quantity correct", "Dented box", "Priority", "Store in a dry place",
	},
	prose: []string{
		"The information in this document is arranged concisely and clearly.",
		"Each section follows a common official document format.",
		"Document numbers and dates are generated consistently from the seed value.",
		"The content can be used to test presentation and layout.",
		"Names and addresses are used for illustration only.",
		"The structure follows common international business correspondence.",
		"The closing section includes space for a signature.",
		"Amounts use US dollar formatting.",
	},
	intro:              "This %s document was generated automatically for sample and demonstration purposes.",
	fictional:          "All fictional data in this document does not refer to any real person, company, or transaction.",
	phonePattern:       "+1-555-###-####",
	addressNumberFirst: true,
	price:              moneyRange{min: 1_000, max: 50_000, div: 100},
	receipt:            moneyRange{min: 10_000, max: 5_000_000, div: 100},
	report:             moneyRange{min: 10_000, max: 1_000_000, div: 100},
}
