package betman

//// getTableBodyInContentsDiv will gets table body in a div
//func getTableBodyInContentsDiv(n *html.Node) []*html.Node {
//	var htmlNode []*html.Node
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		if c.Type != html.ElementNode {
//			continue
//		}
//		if c.Data != crawling.HtmlTagTable {
//			continue
//		}
//
//		if dom.HasClassName(c, "dataH01") {
//			for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
//				if cc.Type != html.ElementNode {
//					continue
//				}
//				if cc.Data != crawling.HtmlTagTableBody {
//					continue
//				}
//				htmlNode = append(htmlNode, cc)
//			}
//		}
//	}
//	return htmlNode
//}
//
//// getTableRowInTableBody will gets table row in a table body
//func getTableRowInTableBody(n *html.Node) []*html.Node {
//	var htmlNode []*html.Node
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		if c.Type != html.ElementNode {
//			continue
//		}
//		if c.Data != crawling.HtmlTagTableRow {
//			continue
//		}
//		htmlNode = append(htmlNode, c)
//	}
//	return htmlNode
//}
//
//// getTableDataInTableRow will gets table data in a table row
//func getTableDataInTableRow(n *html.Node) []*html.Node {
//	var htmlNode []*html.Node
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		if c.Type != html.ElementNode {
//			continue
//		}
//		if c.Data != crawling.HtmlTagTableData {
//			continue
//		}
//		htmlNode = append(htmlNode, c)
//	}
//	return htmlNode
//}
//
//// getGameIdInAnchor will get a href value in a table data
//func getGameIdInAnchor(n *html.Node) string {
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		if c.Type != html.ElementNode {
//			continue
//		}
//		if c.Data != crawling.HtmlTagAnchor {
//			continue
//		}
//		attr, err := dom.GetAttributeByKey(c, "href")
//		if err == nil {
//			return attr.Val
//		}
//	}
//	return ""
//}
