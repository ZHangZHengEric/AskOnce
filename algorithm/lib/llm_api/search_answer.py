from AskOnce.algorithm.lib.llm_api.llm_base_api import LLMBaseAPI

class SearchAnswer (LLMBaseAPI):
    
    def simplify_answer(self,question,search_result,stream=False):
        print('对问题进行简答:',question)
        print('搜索结果长度',len(search_result))
        all_search_result = self.select_search_result(search_result)
        prompt = '''你现在是一个AI智能搜索引擎，根据相关内容（根据问题搜索出来的相关内容），对用户输入的问题进行回答，要求：
1. 回答内容简洁全面，只需要回答问题关注的信息，不要做过多的解释。
2. 只使用相关内容的信息，使用第三人称的表达方法。
3. 不要输出相关信息的来源，例如标题等，也不要输出“根据已知信息”相关的表述。

相关内容：{answers}

问题：{question}'''
        answers = []
        if len(all_search_result)==0:
            print('搜索结果都过滤掉了')
            prompt = '''你现在是一个AI智能助手，对用户输入的问题进行回答，要求：
1. 回答内容简洁全面，只需要回答问题关注的信息，不要做过多的解释。

问题：{question}'''
            prompt_ok = prompt.format(question=question)
        else:
            
            prompt = '''你现在是一个AI智能搜索引擎，根据相关内容（根据问题搜索出来的相关内容），对用户输入的问题进行回答，要求：
1. 回答内容简洁全面，只需要回答问题关注的信息，不要做过多的解释。
2. 只使用相关内容的信息，使用第三人称的表达方法。
3. 不要输出相关信息的来源，例如标题等，也不要输出“根据已知信息”相关的表述。

相关内容：{answers}

问题：{question}'''
            for item in all_search_result:
                answers.append('标题：'+item['title']+'\n部分内容：'+item['answer_for_question'].strip())
            prompt_ok = prompt.format(question=question,answers='\n'.join(answers))
        return_result = {
            'prompt' :prompt_ok,
            'temperature':0.1,
            'presence_penalty':1.2,
            'max_tokens':2048
        }
        if stream:
            return self.ask_llm_stream(**return_result)
        else:
            return self.ask_llm(**return_result)
    
    def detailed_answer(self,question,search_result,stream=False):
        print('对问题进行段落划分的复杂回答:',question)
        print('搜索结果长度',len(search_result))
        all_search_result = self.select_search_result(search_result)
        
        answers = []
        if len(all_search_result)==0:
            print('搜索结果都过滤掉了')
            prompt = '''你现在是一个AI智能助手，对用户输入的问题进行回答，要求：
1. 回答内容丰富详细，字数多，只需要回答问题关注的信息，不要做过多的解释。
2. 使用markdown格式，回答部分进行章节划分，每个段落的文字要内容丰富，字数多。

问题：{question}'''
            prompt_ok = prompt.format(question=question)
        else:
            prompt = '''你现在是一个AI智能搜索引擎，根据相关内容（根据问题搜索出来的结果），对用户输入的问题进行回答，要求：
1. 回答内容丰富详细，字数多，只需要回答问题关注的信息，不要做过多的解释。
2. 只使用相关内容的信息，使用第三人称的表达方法。
3. 不要输出参考资料，不要输出相关信息的来源，例如标题等，也不要输出“根据已知信息”相关的表述。
4. 使用markdown格式，回答部分进行章节划分，每个段落的文字要内容丰富，字数多。

相关内容：{answers}

问题：{question}'''
            for item in all_search_result:
                answers.append('标题：'+item['title']+'\n部分内容：'+item['answer_for_question'].strip())
            prompt_ok = prompt.format(question=question,answers='\n'.join(answers))
        return_result = {
            'prompt' :prompt_ok,
            'temperature':0.1,
            'presence_penalty':1.2,
            'max_tokens':2048
        }
        if stream:
            return self.ask_llm_stream(**return_result)
        else:
            return self.ask_llm(**return_result)
    
    def detailed_no_chapter_answer(self,question,search_result,stream=False):
        print('对问题进行无段落划分的复杂回答:',question)
        print('搜索结果长度',len(search_result))
        all_search_result = self.select_search_result(search_result)
        
        answers = []
        if len(all_search_result)==0:
            print('搜索结果都过滤掉了')
            prompt = '''你现在是一个AI智能助手，根据相关内容（根据问题搜索出来的结果），对用户输入的问题进行回答，要求：
1. 回答内容丰富详细，字数多。
2. 使用markdown格式，但不要有章节层级，使用大的段落。
3. 只需要回答问题关注的信息，不要做过多的解释和内容的延伸。

问题：{question}''' 
            prompt_ok = prompt.format(question=question)
        else:
            prompt = '''你现在是一个AI智能搜索引擎，根据相关内容（根据问题搜索出来的结果），对用户输入的问题进行回答，要求：
1. 回答内容丰富详细，字数多。
2. 只使用相关内容的信息，使用第三人称的表达方法。
3. 不要输出参考资料，不要输出相关信息的来源，例如标题等，也不要输出“根据已知信息”相关的表述。
4. 使用markdown格式，但不要有章节层级，使用大的段落。
5. 只需要回答问题关注的信息，不要做过多的解释和内容的延伸。

相关内容：{answers}

问题：{question}''' 
            for item in all_search_result:
                answers.append('标题：'+item['title']+'\n部分内容：'+item['answer_for_question'].strip())
            prompt_ok = prompt.format(question=question,answers='\n'.join(answers))
        return_result = {
            'prompt' :prompt_ok,
            'temperature':0.1,
            'presence_penalty':1.2,
            'max_tokens':2048
        }
        if stream:
            return self.ask_llm_stream(**return_result)
        else:
            return self.ask_llm(**return_result)
    
    def professional_answer(self,question,search_result,stream=False):
        print('对问题进行专业回答:',question)
        print('搜索结果长度',len(search_result))
        all_search_result = self.select_search_result(search_result)
        prompt = '''你现在是一个AI智能搜索引擎，根据相关内容（根据问题搜索出来的结果），对用户输入的问题进行回答，要求：
1. 回答内容丰富详细，字数多。
2. 只使用相关内容的信息，使用第三人称的表达方法。
3. 不要输出参考资料，不要输出相关信息的来源，例如标题等，也不要输出“根据已知信息”相关的表述。
4. 使用markdown格式，但不要有章节层级，使用大的段落。
5. 只需要回答问题关注的信息，不要做过多的解释和内容的延伸。

相关内容：{answers}

问题：{question}''' 
        answers = []
        if len(all_search_result)==0:
            print('搜索结果都过滤掉了')
            prompt_ok = question
        else:
            for item in all_search_result:
                answers.append('标题：'+item['title']+'\n部分内容：'+item['answer_for_question'].strip())
            prompt_ok = prompt.format(question=question,answers='\n'.join(answers))
        return_result = {
            'prompt' :prompt_ok,
            'temperature':0.1,
            'presence_penalty':1.2,
            'max_tokens':2048
        }
        if stream:
            return self.ask_llm_stream(**return_result)
        else:
            return self.ask_llm(**return_result)
    
    