from AskOnce.algorithm.lib.llm_api.llm_base_api import LLMBaseAPI


class QuestionSplit (LLMBaseAPI):
    def split_question(self,question):
        # 对复杂问题进行拆分
        prompt='''现给你一个主问题，请参考如下示例，当主问题明确询问多个问题时，将主问题拆解成多个子问题，否则输出原始问题。输出要严格按照如下规定：
1. 子问题的每一项使用一行
2. 子问题的每一项是一个关键短句，注意主语明确，不可使用代词来指代前文对象
3. 子问题一定是主问题明确中提到的问题

示例主问题：给出太阳的形状或大小，以及工作原理或发光机制
子问题：
太阳的形状
太阳的大小
太阳的工作原理
太阳的发光机制

示例问题：太阳的发光原理
子问题：
太阳的发光原理

示例问题：斯坦福回应抄袭清华系大模型
子问题：
斯坦福回应抄袭清华系大模型

示例问题：实用水杯推荐
子问题：
实用水杯推荐

主问题：{question}
子问题：
''' 
        quesiton = self.ask_llm(prompt=prompt.format(question=question),temperature=0.2)
        quesiton_list = quesiton.split('\n')
        return quesiton_list