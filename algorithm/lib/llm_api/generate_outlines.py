from AskOnce.algorithm.lib.llm_api.llm_base_api import LLMBaseAPI


class GenerateOutlines (LLMBaseAPI):
    def generate_outlines_by_answer(self,answer):
        prompt='''根据给出的内容，将内容提炼出大纲，大纲的格式要严格按照如下规定，不要有其他的格式：
1. 一级内容使用<h1></h1>包围，不超过15个字，例如<h1>一级内容</h1>
2. 二级内容使用<h2></h2>包围，不超过15个字，例如<h2>二级内容</h2>
3. 三级内容使用<h3></h3>包围，不超过15个字，例如<h3>三级内容</h3>
3. 三级内容使用<h4></h4>包围，不超过15个字，例如<h4>四级内容</h4>

示例内容：中国经济的宏观运行呈现出复杂多样的特征。
首先，从总体经济增长来看，中国国内生产总值（GDP）同比增长5.2%，显示出经济总量的较快增长。这一增长率在全球主要经济体中保持领先，反映出中国经济的强劲动力和全球增长的重要引擎作用。
此外，政府对当前经济形势的判断是“稳”“进”“好”，强调了中国经济回升向好、长期向好的基本趋势没有改变。为了巩固和增强经济回升向好的态势，政府提出了实施扩大内需战略，进一步释放消费潜力，扩大有效投资的工作重点。这些政策措施旨在通过激发消费潜力和扩大投资来促进经济增长。
示例大纲：
<h1>经济增长</h1>
<h2>GDP同比增长5.2%</h2>
<h2>全球主要经济体中增速领先</h2>
<h1>宏观经济政策取向</h1>
<h2>扩大内需战略，释放消费潜力</h2>
<h2>激发有效投资，增强经济回升向好态势</h2>

内容：{answer}
大纲：
'''
        answer_outline = self.ask_llm(prompt=prompt.format(answer=answer),temperature=0.2)
        paper_outline_list = self.parser_outline(answer_outline)
        return paper_outline_list

    def parser_outline(self,paper_outline):    #将大纲文字提取为markdown格式
       # paper_outline='<h1>尼格买提</h1>\n<h2>童年经历</h2>\n<h3>校园欺凌</h3>\n<h2>短视频平台</h2>\n<h3>莫奈花园特效</h3>\n<h2>回应网友</h2>\n<h2>校园欺凌</h2>\n<h2>短视频平台</h2>\n<h3>莫奈花园特效</h3>\n<h2>回应网友</h2>'
        print(paper_outline)
        paper_outline_list  = paper_outline.split('\n')
        paper_outline_out = []
        paper_outline_fir=[]
        paper_outline_sec = []
        paper_outline_third = []
        one_item=paper_outline_list[0]
        for i in range(len(paper_outline_list)):
            one_item=paper_outline_list[i]
            one_item = one_item.strip()
            if one_item.startswith('<h1>'):
                layer = {'level': 'h1', 'content': one_item[len('<h1>'):-len('</h1>')].strip()}
                if len(paper_outline_fir)!=0 and layer['content'] in paper_outline_fir[-1]:
                    paper_outline_sec=None
                    paper_outline_third=None
                    continue
                else:
                    paper_outline_fir.append(layer['content'])
                    paper_outline_out.append(layer)
                    paper_outline_sec=[]
            elif paper_outline_sec!=None and one_item.startswith('<h2>'):
                layer = {'level': 'h2', 'content': one_item[len('<h2>'):-len('</h2>')].strip()}
                if len(paper_outline_sec)!=0 and layer['content'] in paper_outline_sec:
                    paper_outline_third=None
                    continue
                else:
                    paper_outline_sec.append(layer['content'])
                    paper_outline_out.append(layer)
                    paper_outline_third=[]
            elif paper_outline_third!=None and one_item.startswith('<h3>'):
                layer = {'level': 'h3', 'content': one_item[len('<h3>'):-len('</h3>')].strip()}
                if len(paper_outline_third)!=0 and layer['content'] in paper_outline_third:
                    continue
                else:
                    paper_outline_third.append(layer['content'])
                    paper_outline_out.append(layer)
            else:
                continue
        print(paper_outline_out)
        return paper_outline_out